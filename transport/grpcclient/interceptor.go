package grpcclient

import (
	"context"
	"encoding/json"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/constant/sym"
	commonContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const mdMaxLength = bytesize.KB

func newMetadataContextInterceptor(_ context.Context) Interceptor {
	return &metadataContextInterceptor{}
}

func (i *metadataContextInterceptor) unaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = commonContext.MetadataContextToOutgoingContext(ctx)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (i *metadataContextInterceptor) streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ctx = commonContext.MetadataContextToOutgoingContext(ctx)
	return streamer(ctx, desc, cc, method, opts...)
}

func newLoggingInterceptor(_ context.Context, payloadLogSizeLimit uint64) Interceptor {
	return &loggingInterceptor{payloadLogSizeLimit: payloadLogSizeLimit}
}

func (i *loggingInterceptor) unaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fields := make(map[string]interface{})
	fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
	fields[constant.RpcTypeLogKey] = constant.RpcTypeUnary
	fields[constant.MethodLogKey] = method
	fields[constant.IsServerLogKey] = false
	fields[constant.IsRequestLogKey] = true

	// Retrieve outgoing metadata (client-sent)
	outgoingMD, _ := metadata.FromOutgoingContext(ctx)
	truncatedOutgoingMD := utils.TruncateMetadata(outgoingMD, int(mdMaxLength))
	fields[constant.MetadataLogKey] = truncatedOutgoingMD

	// Capture payload
	var message string
	bytes, err := json.Marshal(req)
	if err != nil {
		log.Warn(ctx, err, "error on json.Marshal(req) for GRPC client unaryInterceptor request")
	} else if uint64(len(bytes)) > i.payloadLogSizeLimit {
		message = string(bytes[:i.payloadLogSizeLimit]) + sym.Ellipsis
	} else {
		message = string(bytes)
	}
	fields[constant.MessageLogKey] = message

	// Log to stdout
	log.WithTraceFields(ctx).WithFields(fields).GetLogrusLogger().Info()

	// Variables to capture server-sent metadata
	var incomingMD, trailerMD metadata.MD

	// Add header and trailer options to capture server-sent metadata
	opts = append(opts, grpc.Header(&incomingMD), grpc.Trailer(&trailerMD))

	// Invoke request
	err = invoker(ctx, method, req, reply, cc, opts...)

	fields = make(map[string]interface{})
	fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
	fields[constant.RpcTypeLogKey] = constant.RpcTypeUnary
	fields[constant.MethodLogKey] = method
	fields[constant.IsServerLogKey] = false
	fields[constant.IsRequestLogKey] = false

	// Capture incoming metadata
	truncatedIncomingMD := utils.TruncateMetadata(incomingMD, int(mdMaxLength))
	fields[constant.MetadataLogKey] = truncatedIncomingMD

	// Capture incoming trailers
	truncatedTrailerMD := utils.TruncateMetadata(trailerMD, int(mdMaxLength))
	fields[constant.TrailersLogKey] = truncatedTrailerMD

	if err == nil {
		// Capture payload
		bytes, err := json.Marshal(reply)
		if err != nil {
			log.Warn(ctx, err, "error on json.Marshal(reply) for GRPC client unaryInterceptor response")
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

	return err
}

func (i *loggingInterceptor) streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	cs, err := streamer(ctx, desc, cc, method, opts...)
	return &loggingClientStream{
		ClientStream:        cs,
		method:              method,
		firstSend:           true,
		firstRecv:           true,
		payloadLogSizeLimit: i.payloadLogSizeLimit,
	}, err
}

// SendMsg logs only the first message sent in the stream.
func (s *loggingClientStream) SendMsg(m interface{}) error {
	if s.firstSend {
		ctx := s.Context()
		fields := make(map[string]interface{})
		fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
		fields[constant.RpcTypeLogKey] = constant.RpcTypeStream
		fields[constant.MethodLogKey] = s.method
		fields[constant.IsServerLogKey] = false
		fields[constant.IsRequestLogKey] = true

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
	return s.ClientStream.SendMsg(m)
}

// RecvMsg logs only the first message received in the stream.
func (s *loggingClientStream) RecvMsg(m interface{}) (err error) {
	err = s.ClientStream.RecvMsg(m)
	if s.firstRecv {
		ctx := s.Context()
		fields := make(map[string]interface{})
		fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
		fields[constant.RpcTypeLogKey] = constant.RpcTypeStream
		fields[constant.MethodLogKey] = s.method
		fields[constant.IsServerLogKey] = false
		fields[constant.IsRequestLogKey] = false

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

func (s *loggingClientStream) CloseSend() (err error) {
	// Capture and log trailer metadata when the stream is closed
	trailerMD := s.Trailer()
	err = s.ClientStream.CloseSend()
	if err == nil && len(trailerMD) > 0 {
		ctx := s.Context()
		fields := make(map[string]interface{})
		fields[constant.LogTypeFieldLogKey] = constant.LogTypeGrpc
		fields[constant.RpcTypeLogKey] = constant.RpcTypeStream
		fields[constant.MethodLogKey] = s.method
		fields[constant.IsServerLogKey] = false
		fields[constant.IsRequestLogKey] = false

		truncatedTrailerMD := utils.TruncateMetadata(trailerMD, int(mdMaxLength))
		fields[constant.TrailersLogKey] = truncatedTrailerMD

		// Log to stdout
		log.WithTraceFields(ctx).WithFields(fields).GetLogrusLogger().Info()
	}
	return
}
