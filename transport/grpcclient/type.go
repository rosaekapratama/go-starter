package grpcclient

import (
	"context"
	"google.golang.org/grpc"
)

type IManager interface {
	InitConn(ctx context.Context, connId string, address string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error)
	GetConn(ctx context.Context, connId string) (conn *grpc.ClientConn, err error)
}

type ManagerImpl struct {
	connMap map[string]*grpc.ClientConn
}
