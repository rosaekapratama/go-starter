package restclient

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/headers/contenttype"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/healthcheck"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/transport/constant"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"regexp"
	"time"
)

var (
	_config config.Config
	Manager IManager
)

func Init(ctx context.Context, config config.Config) {
	_config = config
	client, err := newClient(ctx,
		WithLogging(_config.GetObject().Transport.Client.Rest.Logging),
		WithTimeout(time.Duration(_config.GetObject().Transport.Client.Rest.Timeout)*time.Second),
		WithInsecureSkipVerify(_config.GetObject().Transport.Client.Rest.InsecureSkipVerify))
	if err != nil {
		log.Fatal(ctx, err, "Failed to init default rest client")
		return
	}
	Manager = &ManagerImpl{
		defaultClient: client,
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
	if client.logging {
		client.Resty.OnBeforeRequest(preSend)
		client.Resty.OnAfterResponse(postSend)
		client.Resty.OnError(onError)
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

func preSend(_ *resty.Client, r *resty.Request) error {
	if isHealthCheckPath(r.RawRequest) {
		return nil
	}

	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldKey] = constant.LogTypeHttp
	httpFields[constant.UrlFieldKey] = r.URL
	httpFields[constant.MethodFieldKey] = r.Method
	httpFields[constant.IsServerFieldKey] = false
	httpFields[constant.IsRequestFieldKey] = true
	httpFields[constant.HeadersFieldKey] = r.Header
	if r.FormData != nil {
		httpFields[constant.FormDataFieldKey] = r.FormData
	}
	if r.Body != nil {
		httpFields[constant.BodyFieldKey] = r.Body
	}
	log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Info()

	return nil
}

func postSend(_ *resty.Client, r *resty.Response) error {
	if isHealthCheckPath(r.Request.RawRequest) {
		return nil
	}

	httpFields := make(map[string]interface{})
	httpFields[constant.LogTypeFieldKey] = constant.LogTypeHttp
	httpFields[constant.UrlFieldKey] = r.Request.URL
	httpFields[constant.MethodFieldKey] = r.Request.Method
	httpFields[constant.IsServerFieldKey] = false
	httpFields[constant.IsRequestFieldKey] = false
	httpFields[constant.StatusCodeFieldKey] = r.StatusCode()
	httpFields[constant.HeadersFieldKey] = r.Header
	if len(r.Body()) > integer.Zero {
		httpFields[constant.BodyFieldKey] = string(r.Body())
	}
	log.WithTraceFields(r.Request.Context()).WithFields(httpFields).GetLogrusLogger().Info()

	return nil
}

func onError(r *resty.Request, err error) {
	if isHealthCheckPath(r.RawRequest) {
		return
	}

	if v, ok := err.(*resty.ResponseError); ok {
		// v.Response contains the last response from the server
		// v.Err contains the original error

		httpFields := make(map[string]interface{})
		httpFields[constant.LogTypeFieldKey] = constant.LogTypeHttp
		httpFields[constant.UrlFieldKey] = r.URL
		httpFields[constant.MethodFieldKey] = r.Method
		httpFields[constant.IsServerFieldKey] = false
		httpFields[constant.IsRequestFieldKey] = false
		httpFields[constant.ErrorFieldKey] = v.Error()
		log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Error()
	}
	// Log the error, increment a metric, etc...
}

func isHealthCheckPath(r *http.Request) bool {
	isHealthCheck, err := regexp.MatchString(healthcheck.URLPathRegex, r.URL.Path)
	if err != nil {
		log.Errorf(r.Context(), err, "Failed to health path match string, path=%s", r.URL.Path)
		return false
	}
	if isHealthCheck && r.Method == http.MethodGet {
		return true
	}
	return false
}
