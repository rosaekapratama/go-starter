package grpcserver

import (
	"context"
	"encoding/json"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/constant/sym"
	myContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const mdMaxLength = bytesize.KB

func (s *authServerStream) Context() context.Context {
	return s.ctx
}

func (s *loggingServerStream) SendMsg(m interface{}) error {
	if s.firstSend {
		ctx := s.Context()
		fields := make(map[string]interface{})
		fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
		fields[constant.RpcTypeLogKey] = constant.RpcTypeStream
		fields[constant.MethodLogKey] = s.method
		fields[constant.IsServerLogKey] = true
		fields[constant.IsRequestLogKey] = false

		// Retrieve outgoing metadata (client-sent)
		outgoingMD, _ := metadata.FromOutgoingContext(ctx)
		truncatedOutgoingMD := utils.TruncateMetadata(outgoingMD, int(mdMaxLength))
		fields[constant.MetadataLogKey] = truncatedOutgoingMD

		// Capture payload
		var message string
		bytes, err := json.Marshal(m)
		if err != nil {
			log.Warn(ctx, err, "error on json.Marshal(m) for GRPC client streamInterceptor request")
		} else if uint64(len(bytes)) > s.payloadLogSizeLimit {
			message = string(bytes[:s.payloadLogSizeLimit]) + sym.Ellipsis
		} else {
			message = string(bytes)
		}
		fields[constant.MessageLogKey] = message

		// Log to stdout
		log.WithTraceFields(ctx).WithFields(fields).GetLogrusLogger().Info()

		s.firstSend = false // Set the flag to false after logging the first message
	}
	return s.ServerStream.SendMsg(m)
}

func (s *loggingServerStream) RecvMsg(m interface{}) (err error) {
	err = s.ServerStream.RecvMsg(m)
	if s.firstRecv {
		ctx := s.Context()
		fields := make(map[string]interface{})
		fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
		fields[constant.RpcTypeLogKey] = constant.RpcTypeStream
		fields[constant.MethodLogKey] = s.method
		fields[constant.IsServerLogKey] = true
		fields[constant.IsRequestLogKey] = true

		// Capture incoming metadata
		incomingMD, _ := metadata.FromIncomingContext(ctx)
		truncatedIncomingMD := utils.TruncateMetadata(incomingMD, int(mdMaxLength))
		fields[constant.MetadataLogKey] = truncatedIncomingMD

		if err == nil {
			// Capture payload
			var message string
			bytes, err := json.Marshal(m)
			if err != nil {
				log.Warn(ctx, err, "error on json.Marshal(reply) for GRPC client unaryInterceptor response")
			} else if uint64(len(bytes)) > s.payloadLogSizeLimit {
				message = string(bytes[:s.payloadLogSizeLimit]) + sym.Ellipsis
			} else {
				message = string(bytes)
			}
			fields[constant.MessageLogKey] = message
		} else if err != nil {
			fields[constant.ErrorLogKey] = err.Error()
		}
		// Log to stdout
		log.WithTraceFields(ctx).WithFields(fields).GetLogrusLogger().Info()

		s.firstRecv = false // Set the flag to false after logging the first message
	}
	return
}

func newMetadataContextInterceptor(_ context.Context) Interceptor {
	return &metadataContextInterceptor{}
}

func (i *metadataContextInterceptor) unaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx = myContext.MetadataContextFromIncomingContext(ctx)
	return handler(ctx, req)
}

func (i *metadataContextInterceptor) streamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	ctx := myContext.MetadataContextFromIncomingContext(ss.Context())
	return handler(srv, &authServerStream{ss, ctx})
}

func newLoggingInterceptor(_ context.Context, payloadLogSizeLimit uint64) Interceptor {
	return &loggingInterceptor{payloadLogSizeLimit: payloadLogSizeLimit}
}

func (i *loggingInterceptor) unaryInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
	fields := make(map[string]interface{})
	fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
	fields[constant.RpcTypeLogKey] = constant.RpcTypeUnary
	fields[constant.MethodLogKey] = serverInfo.FullMethod
	fields[constant.IsServerLogKey] = true
	fields[constant.IsRequestLogKey] = true

	// Retrieve incoming metadata (client-sent)
	incomingMD, _ := metadata.FromIncomingContext(ctx)
	truncatedIncomingMD := utils.TruncateMetadata(incomingMD, int(mdMaxLength))
	fields[constant.MetadataLogKey] = truncatedIncomingMD

	// Capture payload
	var message string
	bytes, err := json.Marshal(req)
	if err != nil {
		log.Error(ctx, err, "error on json.Marshal(req) for GRPC server unaryInterceptor request")
	} else if uint64(len(bytes)) > i.payloadLogSizeLimit {
		message = string(bytes[:i.payloadLogSizeLimit]) + sym.Ellipsis
	} else {
		message = string(bytes)
	}
	fields[constant.MessageLogKey] = message

	// Log to stdout
	log.WithTraceFields(ctx).WithFields(fields).GetLogrusLogger().Info()

	// Create a new context to capture outgoing response metadata and trailers
	outgoingMD := metadata.MD{}
	trailerMD := metadata.MD{}
	ctx = metadata.NewOutgoingContext(ctx, outgoingMD)

	// Handle request
	res, err = handler(ctx, req)

	fields = make(map[string]interface{})
	fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
	fields[constant.RpcTypeLogKey] = constant.RpcTypeUnary
	fields[constant.MethodLogKey] = serverInfo.FullMethod
	fields[constant.IsServerLogKey] = true
	fields[constant.IsRequestLogKey] = false

	truncatedOutgoingMD := utils.TruncateMetadata(outgoingMD, int(mdMaxLength))
	fields[constant.MetadataLogKey] = truncatedOutgoingMD

	truncatedTrailerMD := utils.TruncateMetadata(trailerMD, int(mdMaxLength))
	fields[constant.TrailersLogKey] = truncatedTrailerMD

	if err == nil {
		// Capture payload
		bytes, err := json.Marshal(res)
		if err != nil {
			log.Error(ctx, err, "error on json.Marshal(reply) for GRPC server unaryInterceptor response")
		} else if uint64(len(bytes)) > i.payloadLogSizeLimit {
			message = string(bytes[:i.payloadLogSizeLimit]) + sym.Ellipsis
		} else {
			message = string(bytes)
		}
		fields[constant.MessageLogKey] = message
	} else {
		fields[constant.ErrorLogKey] = err.Error()
	}

	// Log to stdout
	log.WithTraceFields(ctx).WithFields(fields).GetLogrusLogger().Info()

	return
}

func (i *loggingInterceptor) streamInterceptor(srv any, ss grpc.ServerStream, serverInfo *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	err = handler(srv, &loggingServerStream{
		ServerStream:        ss,
		method:              serverInfo.FullMethod,
		firstSend:           true,
		firstRecv:           true,
		payloadLogSizeLimit: i.payloadLogSizeLimit,
	})
	return
}
