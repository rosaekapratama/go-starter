package elasticsearch

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rosaekapratama/go-starter/transport/restclient"
)

type IManager interface {
	GetClient(ctx context.Context, clientId string) (client *elasticsearch.Client, err error)
}

type managerImpl struct {
	clients map[string]*elasticsearch.Client
}

type loggingTransport struct {
	restClient          *restclient.Client
	isStdoutLogEnabled  bool
	payloadLogSizeLimit int
}
