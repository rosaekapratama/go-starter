package nonce

import (
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/redis"
	"time"
)

const (
	spanSetUpdatePassword = "cache.nonce.SetUpdatePassword"
	spanGetUpdatePassword = "cache.nonce.GetUpdatePassword"

	prefixUpdatePassword = "nonce:password:update:%s:%s"
)

func SetUpdatePassword(ctx context.Context, realm string, userId string, nonce string) error {
	ctx, span := otel.Trace(ctx, spanSetUpdatePassword)
	defer span.End()

	err := redis.Client.Set(ctx, fmt.Sprintf(prefixUpdatePassword, realm, userId), nonce, time.Duration(integer.Three)*time.Minute).Err()
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

func GetUpdatePassword(ctx context.Context, realm string, userId string) (string, error) {
	ctx, span := otel.Trace(ctx, spanGetUpdatePassword)
	defer span.End()

	nonce, err := redis.Client.GetDel(ctx, fmt.Sprintf(prefixUpdatePassword, realm, userId)).Result()
	if err != nil {
		log.Error(ctx, err)
		return str.Empty, err
	}
	return nonce, nil
}
