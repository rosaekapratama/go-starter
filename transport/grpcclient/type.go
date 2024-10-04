package grpcclient

import (
	"context"
	"google.golang.org/grpc"
)

type IManager interface {
	initConn(ctx context.Context, connId string, address string, stdoutLogging bool, databaseLogging string, payloadLogSizeLimit uint64, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error)
	GetConn(ctx context.Context, connId string) (conn *grpc.ClientConn, err error)
}

type managerImpl struct {
	connMap map[string]*grpc.ClientConn
}

type Interceptor interface {
	unaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
	streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error)
}

type metadataContextInterceptor struct {
}

type loggingInterceptor struct {
	payloadLogSizeLimit uint64
}

// loggingClientStream wraps grpc.ClientStream to log stream messages and metadata.
type loggingClientStream struct {
	grpc.ClientStream
	method              string
	firstSend           bool // Flag to indicate the first sent message
	firstRecv           bool // Flag to indicate the first received message
	payloadLogSizeLimit uint64
}
