package grpcserver

import (
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	"google.golang.org/grpc"
	"net"
)

var (
	GRPCServer *grpc.Server
)

func Init(_ context.Context, _ config.Config) {
	var opts []grpc.ServerOption
	GRPCServer = grpc.NewServer(opts...)
}

func Run() {
	ctx := context.Background()
	cfg := config.Instance.GetObject().Transport.Server.Grpc

	// Skip if disabled
	if cfg.Disabled {
		log.Warn(ctx, "GRPC server is disabled")
		return
	}

	port := cfg.Port.Http
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf(ctx, err, "Failed to listen GRPC port, port=%d", port)
	}

	log.Infof(ctx, "Starting GRPC server on port %d", port)
	err = GRPCServer.Serve(lis)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to run GRPC server, port=%d", port)
	}
}
