package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rosaekapratama/go-starter/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
)

var (
	_config config.Config
	client  *Client
)

func Init(config config.Config) {
	_config = config
	client = New()
}

func GetDefaultClient() *Client {
	return client
}

func SetDefaultClient(newClient *Client) {
	client = newClient
}

func New(opts ...ClientOption) *Client {
	cfg := _config.GetObject().Transport.Client.Rest
	client := &Client{Resty: resty.New().SetTransport(
		otelhttp.NewTransport(
			&http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.InsecureSkipVerify,
				},
			},
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("Client HTTP %s %s", r.Method, r.URL.Path)
			}))).
		EnableTrace(),
	}
	for _, opt := range opts {
		opt.Apply(client.Resty)
	}
	return client
}

func NewWithClient(hc *http.Client, opts ...ClientOption) *Client {
	cfg := _config.GetObject().Transport.Client.Rest
	client := &Client{Resty: resty.NewWithClient(hc).SetTransport(
		otelhttp.NewTransport(
			&http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.InsecureSkipVerify,
				},
			},
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("Client HTTP %s %s", r.Method, r.URL.Path)
			}))).
		EnableTrace(),
	}
	for _, opt := range opts {
		opt.Apply(client.Resty)
	}
	return client
}

func NewWithLocalAddr(localAddr net.Addr, opts ...ClientOption) *Client {
	cfg := _config.GetObject().Transport.Client.Rest
	client := &Client{Resty: resty.NewWithLocalAddr(localAddr).SetTransport(
		otelhttp.NewTransport(
			&http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.InsecureSkipVerify,
				},
			},
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("Client HTTP %s %s", r.Method, r.URL.Path)
			}))).
		EnableTrace(),
	}
	for _, opt := range opts {
		opt.Apply(client.Resty)
	}
	return client
}

func newRequest(ctx context.Context, client *Client) *resty.Request {
	return client.Resty.NewRequest().SetContext(ctx)
}

func NewRequest(ctx context.Context) *resty.Request {
	return newRequest(ctx, client)
}

func (c *Client) NewRequest(ctx context.Context) *resty.Request {
	return newRequest(ctx, c)
}
