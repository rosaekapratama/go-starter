package restclient

import (
	"context"
	"crypto/tls"
	"time"
)

type ClientOption interface {
	Apply(ctx context.Context, client *Client) error
}

type loggingOption struct {
	logging bool
}

type insecureSkipVerifyOption struct {
	insecureSkipVerify bool
}

type timeoutClientOption struct {
	timeout time.Duration
}

func (o *loggingOption) Apply(_ context.Context, client *Client) error {
	client.logging = o.logging
	return nil
}

func (o *insecureSkipVerifyOption) Apply(_ context.Context, client *Client) error {
	if o.insecureSkipVerify {
		if client.transport.TLSClientConfig == nil {
			client.transport.TLSClientConfig = &tls.Config{}
		}
		client.transport.TLSClientConfig.InsecureSkipVerify = o.insecureSkipVerify
	}
	return nil
}

func (o *timeoutClientOption) Apply(_ context.Context, client *Client) error {
	client.Resty.SetTimeout(o.timeout)
	return nil
}

func WithLogging(logging bool) ClientOption {
	return &loggingOption{logging: logging}
}

func WithInsecureSkipVerify(insecureSkipVerify bool) ClientOption {
	return &insecureSkipVerifyOption{insecureSkipVerify: insecureSkipVerify}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return &timeoutClientOption{timeout: timeout}
}
