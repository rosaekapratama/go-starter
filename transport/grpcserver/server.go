package grpcserver

import (
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"net"
)

var (
	serverOptions []grpc.ServerOption
	GRPCServer    *grpc.Server
)

func Init(_ context.Context, _ config.Config) {
	opts := make([]grpc.ServerOption, 0)
	opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	opts = append(opts, grpc.UnaryInterceptor(unaryMetadataContextInterceptor))
	opts = append(opts, grpc.StreamInterceptor(streamMetadataContextInterceptor))
	opts = append(opts, serverOptions...)
	GRPCServer = grpc.NewServer(opts...)
}

func AddServerOption(opts ...grpc.ServerOption) {
	serverOptions = append(serverOptions, opts...)
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
