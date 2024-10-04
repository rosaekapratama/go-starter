package grpcserver

import (
	"context"
	"google.golang.org/grpc"
)

// authServerStream is a wrapper around grpc.ServerStream that adds authentication information to the context.
type authServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// loggingServerStream is a wrapper around grpc.ServerStream that adds incoming and outgoing logging
type loggingServerStream struct {
	grpc.ServerStream
	method              string
	firstSend           bool // Flag to indicate the first sent message
	firstRecv           bool // Flag to indicate the first received message
	payloadLogSizeLimit uint64
}

type Interceptor interface {
	unaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
	streamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error)
}

type metadataContextInterceptor struct {
}

type loggingInterceptor struct {
	payloadLogSizeLimit uint64
}
