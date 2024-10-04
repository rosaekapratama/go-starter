package restclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/headers/contenttype"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/log/transport/repositories"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"reflect"
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
	_logDB := _config.GetObject().Transport.Client.Rest.Logging.Database
	_payloadLogSizeLimit, err := bytesize.Parse(_config.GetObject().Transport.Client.Rest.Logging.PayloadLogSizeLimit)
	if err != nil {
		log.Fatal(ctx, err, "Invalid value of REST client payloadLogSizeLimit config")
	}
	payloadLogSizeLimit = int(_payloadLogSizeLimit)

	client, err := newClient(ctx)
	if err != nil {
		log.Fatal(ctx, err, "Failed to init default rest client")
		return
	}
	Manager = &managerImpl{
		defaultClient: client,
	}

	if _logDB != str.Empty {
		LogRepository = restLogRepository
	}
}

func (m *managerImpl) GetDefaultClient() *Client {
	return m.defaultClient
}

func newClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	// List of common config
	logStdout := _config.GetObject().Transport.Client.Rest.Logging.Stdout
	logDB := _config.GetObject().Transport.Client.Rest.Logging.Database
	timeout := time.Duration(_config.GetObject().Transport.Client.Rest.Timeout) * time.Second
	insecureSkipVerify := _config.GetObject().Transport.Client.Rest.Insecure

	client := &Client{
		Resty:     resty.New().EnableTrace(),
		transport: http.DefaultTransport.(*http.Transport),
	}

	// Apply common config
	_ = WithInsecureSkipVerify(insecureSkipVerify).Apply(ctx, client)
	_ = WithTimeout(timeout).Apply(ctx, client)
	_ = WithLogging(logStdout, logDB).Apply(ctx, client)

	// Apply user defined config
	for _, opt := range opts {
		err := opt.Apply(ctx, client)
		if err != nil {
			log.Error(ctx, err, "Failed to apply option to resty client")
			return nil, err
		}
	}

	// Set pre and post of request process
	client.Resty.SetTransport(otelhttp.NewTransport(client.transport))
	if client.logging != nil && client.logging.stdout {
		client.Resty.OnBeforeRequest(preStdoutLogging)
		client.Resty.OnAfterResponse(postStdoutLogging)
		client.Resty.OnError(errorStdoutLogging)
	}
	if client.logging != nil && client.logging.database != str.Empty {
		client.Resty.OnBeforeRequest(preDatabaseLogging)
		client.Resty.OnAfterResponse(postDatabaseLogging)
		client.Resty.OnError(errorDatabaseLogging)
	}
	return client, nil
}

func (m *managerImpl) NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	return newClient(ctx, opts...)
}

func (c *Client) NewRequest(ctx context.Context) *resty.Request {
	request := c.Resty.
		NewRequest().
		SetContext(ctx).
		SetHeader(headers.ContentType, contenttype.ApplicationJson)
	return request
}

func marshalIfContentTypeIsApplicationJson(ctx context.Context, contentType string, body interface{}) (result []byte, ok bool) {
	switch contentType {
	case contenttype.ApplicationJson:
		marshaledBody, err := json.Marshal(body)
		if err != nil {
			t := reflect.TypeOf(body)
			log.Warnf(ctx, "unable to marshal body payload, struct=%s, error=%s", t.Name(), err.Error())
			return
		}
		result = marshaledBody
		ok = true
	}
	return
}

func preStdoutLogging(_ *resty.Client, r *resty.Request) error {
	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeRest
	httpFields[constant.UrlLogKey] = r.URL
	httpFields[constant.MethodLogKey] = r.Method
	httpFields[constant.IsServerLogKey] = false
	httpFields[constant.IsRequestLogKey] = true
	httpFields[constant.HeadersLogKey] = r.Header
	if r.Body != nil {
		var body string
		switch b := r.Body.(type) {
		case string:
			body = b
		case []byte:
			body = string(b)
		default:
			if result, ok := marshalIfContentTypeIsApplicationJson(r.Context(), r.Header.Get(headers.ContentType), b); ok {
				body = string(result)
			} else {
				body = fmt.Sprintf("%v", body)
			}
		}
		if len(body) > payloadLogSizeLimit {
			httpFields[constant.BodyLogKey] = body[:payloadLogSizeLimit] + sym.Ellipsis
		} else if len(body) > 0 {
			httpFields[constant.BodyLogKey] = body
		} else {
			httpFields[constant.BodyLogKey] = str.Empty
		}
	}
	log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Info()
	return nil
}

func postStdoutLogging(_ *resty.Client, r *resty.Response) error {
	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeRest
	httpFields[constant.UrlLogKey] = r.Request.URL
	httpFields[constant.MethodLogKey] = r.Request.Method
	httpFields[constant.IsServerLogKey] = false
	httpFields[constant.IsRequestLogKey] = false
	httpFields[constant.StatusCodeLogKey] = r.StatusCode()
	httpFields[constant.HeadersLogKey] = r.Header
	if len(r.Body()) > payloadLogSizeLimit {
		httpFields[constant.BodyLogKey] = string(r.Body()[:payloadLogSizeLimit]) + sym.Ellipsis
	} else if len(r.Body()) > 0 {
		httpFields[constant.BodyLogKey] = string(r.Body())
	} else {
		httpFields[constant.BodyLogKey] = str.Empty
	}
	log.WithTraceFields(r.Request.Context()).WithFields(httpFields).GetLogrusLogger().Info()

	return nil
}

func errorStdoutLogging(r *resty.Request, err error) {
	var v *resty.ResponseError
	if errors.As(err, &v) {
		// v.Response contains the last response from the server
		// v.Err contains the original error

		httpFields := make(map[string]interface{})
		httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeRest
		httpFields[constant.UrlLogKey] = r.URL
		httpFields[constant.MethodLogKey] = r.Method
		httpFields[constant.IsServerLogKey] = false
		httpFields[constant.IsRequestLogKey] = false
		httpFields[constant.ErrorLogKey] = v.Error()
		log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Error()
	}
	// Log the error, increment a metric, etc...
}

func preDatabaseLogging(_ *resty.Client, r *resty.Request) error {
	LogRepository.SaveRestRequest(r, false)
	return nil
}

func postDatabaseLogging(_ *resty.Client, r *resty.Response) error {
	LogRepository.SaveRestResponse(r, false)
	return nil
}

func errorDatabaseLogging(r *resty.Request, err error) {
	var v *resty.ResponseError
	if errors.As(err, &v) {
		// v.Response contains the last response from the server
		// v.Err contains the original error

		LogRepository.SaveRestError(r, v.Response, false, err)
	}
}
