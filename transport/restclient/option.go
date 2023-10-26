package restclient

import (
	"context"
	"crypto/tls"
	"time"
)

type ClientOption interface {
	Apply(ctx context.Context, client *Client) error
}

type insecureSkipVerifyOption struct {
	insecureSkipVerify bool
}

type timeoutClientOption struct {
	timeout time.Duration
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

func WithInsecureSkipVerify(insecureSkipVerify bool) ClientOption {
	return &insecureSkipVerifyOption{insecureSkipVerify: insecureSkipVerify}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return &timeoutClientOption{timeout: timeout}
}
