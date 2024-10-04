package grpcserver

import (
	"context"
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var (
	serverOptions []grpc.ServerOption
	GRPCServer    *grpc.Server

	defaultPayloadLogSizeLimit = "1KB"
	stdoutLogging              bool
	databaseLogging            string
)

func Init(ctx context.Context, config config.Config) {
	serverConfig := config.GetObject().Transport.Server.Grpc
	payloadLogSizeLimit := defaultPayloadLogSizeLimit
	if serverConfig.Logging != nil {
		stdoutLogging = serverConfig.Logging.Stdout
		if serverConfig.Logging.Database != str.Empty {
			databaseLogging = serverConfig.Logging.Database
		}
		if serverConfig.Logging.PayloadLogSizeLimit != str.Empty {
			payloadLogSizeLimit = serverConfig.Logging.PayloadLogSizeLimit
		}
	}
	_payloadLogSizeLimit, err := bytesize.Parse(payloadLogSizeLimit)
	if err != nil {
		log.Fatal(ctx, err, "Invalid value of GRPC server payloadLogSizeLimit config")
	}

	unaryInterceptorList := make([]grpc.UnaryServerInterceptor, 0)
	streamInterceptorList := make([]grpc.StreamServerInterceptor, 0)
	metadataContextInterceptor := newMetadataContextInterceptor(ctx)
	unaryInterceptorList = append(unaryInterceptorList, metadataContextInterceptor.unaryInterceptor)
	streamInterceptorList = append(streamInterceptorList, metadataContextInterceptor.streamInterceptor)
	if stdoutLogging {
		loggingInterceptor := newLoggingInterceptor(ctx, uint64(_payloadLogSizeLimit))
		unaryInterceptorList = append(unaryInterceptorList, loggingInterceptor.unaryInterceptor)
		streamInterceptorList = append(streamInterceptorList, loggingInterceptor.streamInterceptor)
	}

	opts := make([]grpc.ServerOption, 0)
	opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	opts = append(opts, grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(unaryInterceptorList...)))
	opts = append(opts, grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(streamInterceptorList...)))
	opts = append(opts, serverOptions...)
	GRPCServer = grpc.NewServer(opts...)

	// Enable reflection
	reflection.Register(GRPCServer)
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
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", port))
	if err != nil {
		log.Fatalf(ctx, err, "Failed to listen GRPC port, port=%d", port)
	}

	log.Infof(ctx, "Starting GRPC server on port %d", port)
	err = GRPCServer.Serve(lis)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to run GRPC server, port=%d", port)
	}
}
