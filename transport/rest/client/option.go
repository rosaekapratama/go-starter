package client

import (
	"github.com/go-resty/resty/v2"
	"time"
)

type ClientOption interface {
	Apply(*resty.Client)
}

type timeoutClientOption struct {
	timeout time.Duration
}

func (o *timeoutClientOption) Apply(client *resty.Client) {
	client.SetTimeout(o.timeout)
}
