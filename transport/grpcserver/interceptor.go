package grpcserver

import (
	"context"
	myContext "github.com/rosaekapratama/go-starter/context"
	"google.golang.org/grpc"
)

// authStream is a wrapper around grpc.ServerStream that adds authentication information to the context.
type authStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (a *authStream) Context() context.Context {
	return a.ctx
}

func unaryMetadataContextInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx = myContext.InjectMetadataToContext(ctx)
	return handler(ctx, req)
}

func streamMetadataContextInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	ctx := myContext.InjectMetadataToContext(ss.Context())
	return handler(srv, &authStream{ss, ctx})
}
