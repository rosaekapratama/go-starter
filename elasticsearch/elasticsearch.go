package elasticsearch

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
)

var (
	defaultStdoutLogEnabled    = false
	defaultPayloadLogSizeLimit = 2 * bytesize.KB

	errMissingElasticSearchAddressConfig = "missing elastic search addresses config"

	Manager IManager
)

func Init(ctx context.Context, config config.Config) {
	clients := make(map[string]*elasticsearch.Client)
	esMap := config.GetObject().ElasticSearch
	for clientId, cfg := range esMap {
		if cfg == nil || cfg.Disabled {
			log.Warnf(ctx, "Elastic search client is disabled, instanceId=%s", clientId)
			return
		}

		if cfg.Addresses == nil || len(cfg.Addresses) < integer.One {
			log.Fatal(ctx, response.ConfigNotFound, errMissingElasticSearchAddressConfig)
			return
		}

		// Set default value
		isStdoutLogEnabled := defaultStdoutLogEnabled
		databaseLog := str.Empty
		payloadLogSizeLimit := int(defaultPayloadLogSizeLimit)

		if cfg.Logging != nil {
			isStdoutLogEnabled = cfg.Logging.Stdout
			databaseLog = cfg.Logging.Database
			_payloadLogSizeLimitStr := cfg.Logging.PayloadLogSizeLimit
			if cfg.Logging.PayloadLogSizeLimit != str.Empty {
				_payloadLogSizeLimit, err := bytesize.Parse(_payloadLogSizeLimitStr)
				if err != nil {
					log.Fatalf(ctx, err, "error on bytesize.Parse() of payload log size limit, clientId=%s, payloadLogSizeLimit=%s", clientId, _payloadLogSizeLimitStr)
				}
				payloadLogSizeLimit = int(_payloadLogSizeLimit)
			}
		}

		var err error
		client, err := elasticsearch.NewClient(
			elasticsearch.Config{
				Addresses: cfg.Addresses,
				Username:  cfg.Username,
				Password:  cfg.Password,
				Transport: NewLoggingTransport(ctx, isStdoutLogEnabled, databaseLog, payloadLogSizeLimit),
			})
		if err != nil {
			log.Fatal(ctx, err, "Failed to create elastic search client")
			return
		}
		clients[clientId] = client
		log.Infof(ctx, "Elastic search is initiated, instanceId=%s, address=%v", clientId, cfg.Addresses)
	}
	Manager = &managerImpl{clients: clients}
}

func (m *managerImpl) GetClient(ctx context.Context, clientId string) (client *elasticsearch.Client, err error) {
	var exists bool
	if client, exists = m.clients[clientId]; !exists {
		log.Errorf(ctx, response.GeneralError, "unregistered elastic search client, clientId=%s", clientId)
		err = response.GeneralError
	}
	return
}
