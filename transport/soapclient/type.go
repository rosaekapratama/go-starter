package soapclient

import (
	"context"
	"github.com/tiaguinho/gosoap"
	"net/http"
)

type IManager interface {
	GetDefaultClient() *Client
	NewClient(ctx context.Context, opts ...ClientOption) (*Client, error)
}

type ManagerImpl struct {
	defaultClient *Client
}

type Client struct {
	httpClient    *http.Client
	transport     *loggingTransport
	soapClientMap map[string]*gosoap.Client
}

type clientLogging struct {
	Stdout   bool   `yaml:"stdout"`
	Database string `yaml:"database"`
}

type loggingTransport struct {
	transport *http.Transport
	logging   *clientLogging
}
