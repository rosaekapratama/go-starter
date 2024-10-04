package grpcclient

import (
	"context"
	"errors"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	defaultPayloadLogSizeLimit = "1KB"
	errGRPCClientConnNotFound  = errors.New("GRPC client connection not found")

	Manager IManager
)

func Init(ctx context.Context, config config.Config) {
	Manager = &managerImpl{connMap: make(map[string]*grpc.ClientConn)}
	connMap := config.GetObject().Transport.Client.Grpc
	for connId, connConfig := range connMap {
		var stdoutLogging bool
		var databaseLogging string
		payloadLogSizeLimit := defaultPayloadLogSizeLimit
		if connConfig.Logging != nil {
			stdoutLogging = connConfig.Logging.Stdout
			if connConfig.Logging.Database != str.Empty {
				databaseLogging = connConfig.Logging.Database
			}
			if connConfig.Logging.PayloadLogSizeLimit != str.Empty {
				payloadLogSizeLimit = connConfig.Logging.PayloadLogSizeLimit
			}
		}
		_payloadLogSizeLimit, err := bytesize.Parse(payloadLogSizeLimit)
		if err != nil {
			log.Fatal(ctx, err, "Invalid value of GRPC client payloadLogSizeLimit config")
		}

		opts := make([]grpc.DialOption, 0)
		if connConfig.Insecure {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
		_, err = Manager.initConn(ctx, connId, connConfig.Address, stdoutLogging, databaseLogging, uint64(_payloadLogSizeLimit), opts...)
		if err != nil {
			log.Fatalf(ctx, err, "error on init GRPC connection, connId=%s", connId)
		} else {

		}
	}
}

func (m managerImpl) initConn(ctx context.Context, connId string, address string, stdoutLogging bool, databaseLogging string, payloadLogSizeLimit uint64, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	grpcOptions := append([]grpc.DialOption{}, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	metadataContextInterceptor := newMetadataContextInterceptor(ctx)
	grpcOptions = append(grpcOptions, grpc.WithUnaryInterceptor(metadataContextInterceptor.unaryInterceptor))
	grpcOptions = append(grpcOptions, grpc.WithStreamInterceptor(metadataContextInterceptor.streamInterceptor))
	if stdoutLogging {
		loggingInterceptor := newLoggingInterceptor(ctx, payloadLogSizeLimit)
		grpcOptions = append(grpcOptions, grpc.WithUnaryInterceptor(loggingInterceptor.unaryInterceptor))
		grpcOptions = append(grpcOptions, grpc.WithStreamInterceptor(loggingInterceptor.streamInterceptor))
	}
	grpcOptions = append(grpcOptions, opts...)
	conn, err = grpc.Dial(address, grpcOptions...)
	if err != nil {
		log.Errorf(ctx, err, "failed to init GRPC client conn, connId=%s", connId)
		return
	}
	m.connMap[connId] = conn
	log.Infof(ctx, "GRPC client connection is initiated, connId=%s, address=%s", connId, address)
	return
}

func (m managerImpl) GetConn(ctx context.Context, connId string) (conn *grpc.ClientConn, err error) {
	if v, exists := m.connMap[connId]; exists {
		conn = v
		return
	}

	err = errGRPCClientConnNotFound
	log.Errorf(ctx, err, "GRPC client connection not found, connId=%s", connId)
	return
}
