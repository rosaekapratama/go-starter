package soapclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/headers/contenttype"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/log/transport/repositories"
	"github.com/rosaekapratama/go-starter/utils"
	"github.com/tiaguinho/gosoap"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"time"
)

var (
	_config             config.Config
	Manager             IManager
	LogRepository       repositories.ITransportLogRepository
	payloadLogSizeLimit int
)

func Init(ctx context.Context, config config.Config, restLogRepository repositories.ITransportLogRepository) {
	_config = config
	logStdout := _config.GetObject().Transport.Client.Soap.Logging.Stdout
	logDB := _config.GetObject().Transport.Client.Soap.Logging.Database
	_payloadLogSizeLimit, err := bytesize.Parse(_config.GetObject().Transport.Client.Soap.Logging.PayloadLogSizeLimit)
	if err != nil {
		log.Fatal(ctx, err, "Invalid value of SOAP client payloadLogSizeLimit config")
	}
	payloadLogSizeLimit = int(_payloadLogSizeLimit)

	client, err := newClient(ctx,
		WithLogging(logStdout, logDB),
		WithTimeout(time.Duration(_config.GetObject().Transport.Client.Rest.Timeout)*time.Second),
		WithInsecureSkipVerify(_config.GetObject().Transport.Client.Rest.Insecure))
	if err != nil {
		log.Fatal(ctx, err, "Failed to init default rest client")
		return
	}
	Manager = &managerImpl{
		defaultClient: client,
	}

	if logDB != str.Empty {
		LogRepository = restLogRepository
	}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request details
	var clonedReq *http.Request
	var clearFunc func()
	var err error
	if t.logging != nil {
		clonedReq, clearFunc, err = utils.CloneHttpRequest(req, payloadLogSizeLimit)
		defer clearFunc()
		if err != nil {
			log.Error(req.Context(), err, "failed to clone http request for soap logging")
		}

		if t.logging.Stdout {
			go preStdoutLogging(clonedReq)
		}
		if t.logging.Database != str.Empty {
			go preDatabaseLogging(clonedReq)
		}
	}

	// Perform the actual HTTP request
	res, err := t.transport.RoundTrip(req)
	if err != nil && t.logging != nil {
		if t.logging.Stdout {
			go errorStdoutLogging(clonedReq, err)
		}
		if t.logging.Database != str.Empty {
			go errorDatabaseLogging(clonedReq, err)
		}
	}

	// Log the response details
	if res != nil && t.logging != nil {
		clonedRes, clearFunc, err := utils.CloneHttpResponse(res, payloadLogSizeLimit)
		defer clearFunc()
		if err != nil {
			log.Error(req.Context(), err, "failed to clone http response for soap logging")
		}

		if t.logging.Stdout {
			go postStdoutLogging(clonedRes)
		}
		if t.logging.Database != str.Empty {
			go postDatabaseLogging(clonedRes)
		}
	}

	return res, err
}

func (m *managerImpl) GetDefaultClient() *Client {
	return m.defaultClient
}

func newClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	transport := &loggingTransport{transport: http.DefaultTransport.(*http.Transport)}
	client := &Client{
		httpClient: &http.Client{Transport: otelhttp.NewTransport(transport)},
		transport:  transport,
	}
	for _, opt := range opts {
		err := opt.Apply(ctx, client)
		if err != nil {
			log.Error(ctx, err, "Failed to apply option to http client")
			return nil, err
		}
	}
	return client, nil
}

func (m *managerImpl) NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	return newClient(ctx, opts...)
}

func (c *Client) Call(ctx context.Context, wsdlAddress string, method string, params gosoap.SoapParams) (res *gosoap.Response, err error) {
	var soapClient *gosoap.Client
	var exists bool
	if soapClient, exists = c.soapClientMap[wsdlAddress]; !exists {
		soapClient, err = gosoap.SoapClient(wsdlAddress, c.httpClient)
		if err != nil {
			log.Errorf(ctx, err, "failed to get SOAP client, wsdlAddress=%s", wsdlAddress)
			return
		}
		c.soapClientMap[wsdlAddress] = soapClient
	}

	res, err = soapClient.Call(method, params)
	if err != nil {
		log.Error(ctx, err, "failed to call SOAP method, wsdlAddress=%s, method=%s", wsdlAddress, method)
	}
	return
}

func preStdoutLogging(r *http.Request) {
	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeSoap
	httpFields[constant.UrlLogKey] = r.URL
	httpFields[constant.SoapActionLogKey] = r.Header.Get(headers.SOAPAction)
	httpFields[constant.MethodLogKey] = r.Method
	httpFields[constant.IsServerLogKey] = false
	httpFields[constant.IsRequestLogKey] = true
	httpFields[constant.HeadersLogKey] = r.Header
	if r.Body != nil {
		var body string
		contentType := r.Header.Get(headers.ContentType)
		if contentType == contenttype.ApplicationJson {
			bytes, err := json.Marshal(r.Body)
			if err != nil {
				log.Error(r.Context(), err, "Failed to marshal request body payload for logging")
				body = "Cannot parsed payload"
			}
			body = string(bytes)
		} else {
			body = fmt.Sprintf("%v", r.Body)
		}
		if len(body) > payloadLogSizeLimit {
			httpFields[constant.BodyLogKey] = body[:payloadLogSizeLimit]
		} else if len(body) > 0 {
			httpFields[constant.BodyLogKey] = body
		} else {
			httpFields[constant.BodyLogKey] = str.Empty
		}
	}
	log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Info()
}

func postStdoutLogging(r *http.Response) {
	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeSoap
	httpFields[constant.UrlLogKey] = r.Request.URL
	httpFields[constant.SoapActionLogKey] = r.Request.Header.Get(headers.SOAPAction)
	httpFields[constant.MethodLogKey] = r.Request.Method
	httpFields[constant.IsServerLogKey] = false
	httpFields[constant.IsRequestLogKey] = false
	httpFields[constant.StatusCodeLogKey] = r.StatusCode
	httpFields[constant.HeadersLogKey] = r.Header
	if r.Body != nil {
		var body string
		contentType := r.Header.Get(headers.ContentType)
		if contentType == contenttype.ApplicationJson {
			bytes, err := json.Marshal(r.Body)
			if err != nil {
				log.Error(r.Request.Context(), err, "Failed to marshal response body payload for logging")
				body = "Cannot parsed payload"
			}
			body = string(bytes)
		} else {
			body = fmt.Sprintf("%v", r.Body)
		}
		if len(body) > payloadLogSizeLimit {
			httpFields[constant.BodyLogKey] = body[:payloadLogSizeLimit]
		} else if len(body) > 0 {
			httpFields[constant.BodyLogKey] = body
		} else {
			httpFields[constant.BodyLogKey] = str.Empty
		}
	}
	log.WithTraceFields(r.Request.Context()).WithFields(httpFields).GetLogrusLogger().Info()
}

func errorStdoutLogging(r *http.Request, err error) {
	if v, ok := err.(*resty.ResponseError); ok {
		// v.Response contains the last response from the server
		// v.Err contains the original error

		httpFields := make(map[string]interface{})
		httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeRest
		httpFields[constant.UrlLogKey] = r.URL
		httpFields[constant.SoapActionLogKey] = r.Header.Get(headers.SOAPAction)
		httpFields[constant.MethodLogKey] = r.Method
		httpFields[constant.IsServerLogKey] = false
		httpFields[constant.IsRequestLogKey] = false
		httpFields[constant.ErrorLogKey] = v.Error()
		log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Error()
	}
	// Log the error, increment a metric, etc...
}

func preDatabaseLogging(r *http.Request) {
	LogRepository.SaveSoapRequest(r, r.Header.Get(headers.SOAPAction), false)
}

func postDatabaseLogging(r *http.Response) {
	LogRepository.SaveSoapResponse(r, r.Request.Header.Get(headers.SOAPAction), false)
}

func errorDatabaseLogging(r *http.Request, err error) {
	LogRepository.SaveSoapError(r, r.Header.Get(headers.SOAPAction), false, err)
}
