package soapclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/headers/contenttype"
	"github.com/rosaekapratama/go-starter/constant/integer"
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
	_config       config.Config
	Manager       IManager
	LogRepository repositories.ITransportLogRepository
)

func Init(ctx context.Context, config config.Config, restLogRepository repositories.ITransportLogRepository) {
	_config = config
	logStdout := _config.GetObject().Transport.Client.Soap.Logging.Stdout
	logDB := _config.GetObject().Transport.Client.Soap.Logging.Database
	client, err := newClient(ctx,
		WithLogging(logStdout, logDB),
		WithTimeout(time.Duration(_config.GetObject().Transport.Client.Rest.Timeout)*time.Second),
		WithInsecureSkipVerify(_config.GetObject().Transport.Client.Rest.InsecureSkipVerify))
	if err != nil {
		log.Fatal(ctx, err, "Failed to init default rest client")
		return
	}
	Manager = &ManagerImpl{
		defaultClient: client,
	}

	if logDB != str.Empty {
		LogRepository = restLogRepository
	}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request details
	if t.logging != nil {
		r, err := utils.CloneHttpRequest(req)
		if err != nil {
			log.Error(req.Context(), err, "failed to clone http request for soap logging")
		}

		if t.logging.Stdout {
			go preStdoutLogging(r)
		}
		if t.logging.Database != str.Empty {
			go preDatabaseLogging(r)
		}
	}

	// Perform the actual HTTP request
	res, err := t.transport.RoundTrip(req)
	if err != nil && t.logging != nil {
		r, errClone := utils.CloneHttpRequest(req)
		if errClone != nil {
			log.Error(req.Context(), errClone, "failed to clone http error for soap logging")
		}

		if t.logging.Stdout {
			go errorStdoutLogging(r, err)
		}
		if t.logging.Database != str.Empty {
			go errorDatabaseLogging(r, err)
		}
	}

	// Log the response details
	if res != nil && t.logging != nil {
		r, err := utils.CloneHttpResponse(res)
		if err != nil {
			log.Error(req.Context(), err, "failed to clone http response for soap logging")
		}

		if t.logging.Stdout {
			go postStdoutLogging(r)
		}
		if t.logging.Database != str.Empty {
			go postDatabaseLogging(r)
		}
	}

	return res, err
}

func (m *ManagerImpl) GetDefaultClient() *Client {
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

func (m *ManagerImpl) NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
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
	httpFields[constant.LogTypeFieldKey] = constant.LogTypeSoap
	httpFields[constant.UrlFieldKey] = r.URL
	httpFields[constant.SoapActionFieldKey] = r.Header.Get(headers.SOAPAction)
	httpFields[constant.MethodFieldKey] = r.Method
	httpFields[constant.IsServerFieldKey] = false
	httpFields[constant.IsRequestFieldKey] = true
	httpFields[constant.HeadersFieldKey] = r.Header
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
		if len(body) > integer.Zero && len(body) <= (64*1000) {
			httpFields[constant.BodyFieldKey] = body
		}
	}
	log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Info()
}

func postStdoutLogging(r *http.Response) {
	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldKey] = constant.LogTypeSoap
	httpFields[constant.UrlFieldKey] = r.Request.URL
	httpFields[constant.SoapActionFieldKey] = r.Request.Header.Get(headers.SOAPAction)
	httpFields[constant.MethodFieldKey] = r.Request.Method
	httpFields[constant.IsServerFieldKey] = false
	httpFields[constant.IsRequestFieldKey] = false
	httpFields[constant.StatusCodeFieldKey] = r.StatusCode
	httpFields[constant.HeadersFieldKey] = r.Header
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
		if len(body) > integer.Zero && len(body) <= (64*1000) {
			httpFields[constant.BodyFieldKey] = body
		}
	}
	log.WithTraceFields(r.Request.Context()).WithFields(httpFields).GetLogrusLogger().Info()
}

func errorStdoutLogging(r *http.Request, err error) {
	if v, ok := err.(*resty.ResponseError); ok {
		// v.Response contains the last response from the server
		// v.Err contains the original error

		httpFields := make(map[string]interface{})
		httpFields[constant.LogTypeFieldKey] = constant.LogTypeRest
		httpFields[constant.UrlFieldKey] = r.URL
		httpFields[constant.SoapActionFieldKey] = r.Header.Get(headers.SOAPAction)
		httpFields[constant.MethodFieldKey] = r.Method
		httpFields[constant.IsServerFieldKey] = false
		httpFields[constant.IsRequestFieldKey] = false
		httpFields[constant.ErrorFieldKey] = v.Error()
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
