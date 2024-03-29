package restclient

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
	logStdout := _config.GetObject().Transport.Client.Rest.Logging.Stdout
	logDB := _config.GetObject().Transport.Client.Rest.Logging.Database
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

func (m *ManagerImpl) GetDefaultClient() *Client {
	return m.defaultClient
}

func newClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	client := &Client{
		Resty:     resty.New().EnableTrace(),
		transport: http.DefaultTransport.(*http.Transport),
	}
	for _, opt := range opts {
		err := opt.Apply(ctx, client)
		if err != nil {
			log.Error(ctx, err, "Failed to apply option to resty client")
			return nil, err
		}
	}
	client.Resty.SetTransport(otelhttp.NewTransport(client.transport))
	if client.logging != nil && client.logging.Stdout {
		client.Resty.OnBeforeRequest(preStdoutLogging)
		client.Resty.OnAfterResponse(postStdoutLogging)
		client.Resty.OnError(errorStdoutLogging)
	}
	if client.logging != nil && client.logging.Database != str.Empty {
		client.Resty.OnBeforeRequest(preDatabaseLogging)
		client.Resty.OnAfterResponse(postDatabaseLogging)
		client.Resty.OnError(errorDatabaseLogging)
	}
	return client, nil
}

func (m *ManagerImpl) NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	return newClient(ctx, opts...)
}

func (c *Client) NewRequest(ctx context.Context) *resty.Request {
	request := c.Resty.
		NewRequest().
		SetContext(ctx).
		SetHeader(headers.ContentType, contenttype.ApplicationJson)
	return request
}

func preStdoutLogging(_ *resty.Client, r *resty.Request) error {
	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldKey] = constant.LogTypeRest
	httpFields[constant.UrlFieldKey] = r.URL
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
				log.Error(r.Context(), err, "Failed to marshal body payload for logging")
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
	return nil
}

func postStdoutLogging(_ *resty.Client, r *resty.Response) error {
	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldKey] = constant.LogTypeRest
	httpFields[constant.UrlFieldKey] = r.Request.URL
	httpFields[constant.MethodFieldKey] = r.Request.Method
	httpFields[constant.IsServerFieldKey] = false
	httpFields[constant.IsRequestFieldKey] = false
	httpFields[constant.StatusCodeFieldKey] = r.StatusCode()
	httpFields[constant.HeadersFieldKey] = r.Header
	if len(r.Body()) > integer.Zero && len(r.Body()) <= (64*1000) {
		httpFields[constant.BodyFieldKey] = string(r.Body())
	}
	log.WithTraceFields(r.Request.Context()).WithFields(httpFields).GetLogrusLogger().Info()

	return nil
}

func errorStdoutLogging(r *resty.Request, err error) {
	if v, ok := err.(*resty.ResponseError); ok {
		// v.Response contains the last response from the server
		// v.Err contains the original error

		httpFields := make(map[string]interface{})
		httpFields[constant.LogTypeFieldKey] = constant.LogTypeRest
		httpFields[constant.UrlFieldKey] = r.URL
		httpFields[constant.MethodFieldKey] = r.Method
		httpFields[constant.IsServerFieldKey] = false
		httpFields[constant.IsRequestFieldKey] = false
		httpFields[constant.ErrorFieldKey] = v.Error()
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
	if v, ok := err.(*resty.ResponseError); ok {
		// v.Response contains the last response from the server
		// v.Err contains the original error

		LogRepository.SaveRestError(r, v.Response, false, err)
	} else {
		LogRepository.SaveRestError(r, nil, false, err)
	}
}
