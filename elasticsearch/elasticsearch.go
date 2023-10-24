package elasticsearch

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
)

var (
	errMissingElasticSearchAddressConfig = "missing elastic search addresses config"

	Client *elasticsearch.Client
)

func Init(ctx context.Context, config config.Config) {
	cfg := config.GetObject().ElasticSearch
	if cfg == nil || cfg.Disabled {
		log.Warn(ctx, "Elastic search client is disabled")
		return
	}

	if cfg.Addresses == nil || len(cfg.Addresses) < integer.One {
		log.Fatal(ctx, response.ConfigNotFound, errMissingElasticSearchAddressConfig)
		return
	}

	var err error
	Client, err = elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: cfg.Addresses,
			Username:  cfg.Username,
			Password:  cfg.Password,
		})
	if err != nil {
		log.Fatal(ctx, err, "Failed to create elastic search client")
		return
	}
}
