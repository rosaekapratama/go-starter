package sftpclient

import (
	"context"
	"errors"
	"github.com/rosaekapratama/go-starter/config"
)

var (
	errSftpClientConnNotFound = errors.New("GRPC client connection not found")

	Manager IManager
)

func Init(ctx context.Context, config config.Config) {

}
