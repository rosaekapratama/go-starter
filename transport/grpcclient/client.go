package grpcclient

import (
	"context"
	"errors"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var (
	errGRPCClientConnNotFound = errors.New("GRPC client connection not found")

	Manager IManager
)

func Init(_ context.Context, _ config.Config) {
	Manager = &ManagerImpl{connMap: make(map[string]*grpc.ClientConn)}
}

func (m ManagerImpl) InitConn(ctx context.Context, connId string, address string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	grpcOptions := append([]grpc.DialOption{}, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	grpcOptions = append(grpcOptions, grpc.WithUnaryInterceptor(unaryMetadataContextInterceptor))
	grpcOptions = append(grpcOptions, grpc.WithStreamInterceptor(streamMetadataContextInterceptor))
	grpcOptions = append(grpcOptions, opts...)
	conn, err = grpc.Dial(address, grpcOptions...)
	if err != nil {
		log.Errorf(ctx, err, "failed to init GRPC client conn, connId=%s", connId)
		return
	}
	m.connMap[connId] = conn
	return
}

func (m ManagerImpl) GetConn(ctx context.Context, connId string) (conn *grpc.ClientConn, err error) {
	if v, exists := m.connMap[connId]; exists {
		conn = v
		return
	}

	err = errGRPCClientConnNotFound
	log.Errorf(ctx, err, "GRPC client connection is not found, connId=%s", connId)
	return
}
