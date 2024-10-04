package response

import (
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
)

type IResponse interface {
	Code() string
	Description() string
	HttpStatusCode() int
	OtelCode() otelCodes.Code
	GrpcCode() grpcCodes.Code
	IsError() bool
}

type Response int
