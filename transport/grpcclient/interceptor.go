package grpcclient

import (
	"context"
	myContext "github.com/rosaekapratama/go-starter/context"
	"google.golang.org/grpc"
)

func unaryMetadataContextInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = myContext.InjectContextToMetadata(ctx)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func streamMetadataContextInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ctx = myContext.InjectContextToMetadata(ctx)
	return streamer(ctx, desc, cc, method, opts...)
}
