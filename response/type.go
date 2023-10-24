package response

import "go.opentelemetry.io/otel/codes"

type IResponse interface {
	Code() string
	Description() string
	HttpStatusCode() int
	OtelCode() codes.Code
}

type Response int
