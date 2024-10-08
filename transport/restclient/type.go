package restclient

import (
	"context"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type IManager interface {
	GetDefaultClient() *Client
	NewClient(ctx context.Context, opts ...ClientOption) (*Client, error)
}

type managerImpl struct {
	defaultClient *Client
}

type Client struct {
	Resty     *resty.Client
	transport *http.Transport
	logging   *clientLogging
}

type clientLogging struct {
	stdout   bool
	database string
}
