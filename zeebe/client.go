package zeebe

import (
	"context"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
	"os"
)

const (
	zeebeAddressEnvVarKey                = "ZEEBE_ADDRESS"
	zeebeClientIdEnvVarKey               = "ZEEBE_CLIENT_ID"
	zeebeClientSecretEnvVarKey           = "ZEEBE_CLIENT_SECRET"
	zeebeAuthorizationServerUrlEnvVarKey = "ZEEBE_AUTHORIZATION_SERVER_URL"
)

var (
	Client zbc.Client
)

func Init(ctx context.Context, config config.Config) {
	cfg := config.GetObject().Zeebe

	if cfg == nil || cfg.Disabled {
		log.Warn(ctx, "Zeebe client is disabled")
		return
	}

	// validate config
	if cfg.Address == str.Empty {
		log.Fatalf(ctx, response.ConfigNotFound, "Zeebe address is missing")
		return
	}
	if cfg.ClientId == str.Empty {
		log.Fatalf(ctx, response.ConfigNotFound, "Zeebe client ID is missing")
		return
	}
	if cfg.ClientSecret == str.Empty {
		log.Fatalf(ctx, response.ConfigNotFound, "Zeebe client secret is missing")
		return
	}

	var err error
	err = os.Setenv(zeebeAddressEnvVarKey, cfg.Address)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to set env var, key=%s, value=%s", zeebeAddressEnvVarKey, cfg.Address)
		return
	}
	err = os.Setenv(zeebeClientIdEnvVarKey, cfg.ClientId)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to set env var, key=%s, value=%s", zeebeClientIdEnvVarKey, cfg.ClientId)
		return
	}
	err = os.Setenv(zeebeClientSecretEnvVarKey, cfg.ClientSecret)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to set env var, key=%s, value=%s", zeebeClientSecretEnvVarKey, cfg.ClientSecret)
		return
	}
	err = os.Setenv(zeebeAuthorizationServerUrlEnvVarKey, cfg.AuthorizationServerURL)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to set env var, key=%s, value=%s", zeebeAuthorizationServerUrlEnvVarKey, cfg.AuthorizationServerURL)
		return
	}

	Client, err = zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress: cfg.Address,
	})
	if err != nil {
		log.Fatalf(ctx, err, "Failed to initiate zeebe client")
		return
	}

	res, err := Client.NewTopologyCommand().Send(ctx)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to get topology response")
		return
	}
	log.Infof(ctx, "Zeebe client is initiated, gatewayVersion=%s, clusterSize=%d", res.GetGatewayVersion(), res.GetClusterSize())
}
