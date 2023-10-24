package response

import (
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

const (
	Success Response = iota
	GeneralError
	InitFailed
	InvalidConfigValueType
	InvalidArgument
	InvalidPageRequest
	DataNotFound
	ConfigNotFound
	UnknownResponse
	APINotRegistered
	DataIsEmpty
	UnauthorizedAccess
	InvalidBodyRequest
	OperationNotPermitted
	RequestEntityTooLarge
	InvalidCredentials
	AvroSchemaNotFound
	InvalidConfig
	Processing
	UnsupportedType
	MissingGoogleToken
)

var (
	descriptions = map[Response]string{
		Success:                "Success",
		GeneralError:           "General error",
		InitFailed:             "Init failed",
		InvalidConfigValueType: "Invalid config value type",
		InvalidArgument:        "Invalid argument",
		InvalidPageRequest:     "Invalid page request",
		DataNotFound:           "Data not found",
		ConfigNotFound:         "Config not found",
		UnknownResponse:        "Unknown response",
		APINotRegistered:       "API is not registered, path=%s, method=%s",
		DataIsEmpty:            "Data is empty",
		UnauthorizedAccess:     "Unauthorized access",
		InvalidBodyRequest:     "Invalid body request",
		OperationNotPermitted:  "Operation not permitted",
		RequestEntityTooLarge:  "Request entity too large",
		InvalidCredentials:     "Invalid credentials",
		AvroSchemaNotFound:     "Avro schema not found",
		InvalidConfig:          "Invalid config",
		Processing:             "Processing",
		UnsupportedType:        "Unsupported type",
		MissingGoogleToken:     "Missing google token",
	}

	httpStatusCodes = map[Response]int{
		Success:                http.StatusOK,
		GeneralError:           http.StatusInternalServerError,
		InitFailed:             http.StatusInternalServerError,
		InvalidConfigValueType: http.StatusInternalServerError,
		InvalidArgument:        http.StatusBadRequest,
		InvalidPageRequest:     http.StatusBadRequest,
		DataNotFound:           http.StatusOK,
		ConfigNotFound:         http.StatusInternalServerError,
		UnknownResponse:        http.StatusInternalServerError,
		APINotRegistered:       http.StatusNotFound,
		DataIsEmpty:            http.StatusOK,
		UnauthorizedAccess:     http.StatusUnauthorized,
		InvalidBodyRequest:     http.StatusBadRequest,
		OperationNotPermitted:  http.StatusForbidden,
		RequestEntityTooLarge:  http.StatusRequestEntityTooLarge,
		InvalidCredentials:     http.StatusUnauthorized,
		AvroSchemaNotFound:     http.StatusInternalServerError,
		InvalidConfig:          http.StatusInternalServerError,
		Processing:             http.StatusProcessing,
		UnsupportedType:        http.StatusInternalServerError,
		MissingGoogleToken:     http.StatusBadRequest,
	}

	otelCodes = map[Response]codes.Code{
		Success:                codes.Ok,
		GeneralError:           codes.Error,
		InitFailed:             codes.Error,
		InvalidConfigValueType: codes.Error,
		InvalidArgument:        codes.Ok,
		InvalidPageRequest:     codes.Ok,
		DataNotFound:           codes.Ok,
		UnknownResponse:        codes.Error,
		APINotRegistered:       codes.Ok,
		DataIsEmpty:            codes.Ok,
		UnauthorizedAccess:     codes.Ok,
		InvalidBodyRequest:     codes.Ok,
		OperationNotPermitted:  codes.Ok,
		RequestEntityTooLarge:  codes.Ok,
		InvalidCredentials:     codes.Ok,
		AvroSchemaNotFound:     codes.Ok,
		InvalidConfig:          codes.Ok,
		Processing:             codes.Ok,
		UnsupportedType:        codes.Ok,
		MissingGoogleToken:     codes.Ok,
	}
)

func (r Response) Code() string {
	return fmt.Sprintf("%04d", r)
}

func (r Response) Description() string {
	return descriptions[r]
}

func (r Response) Error() string {
	return r.Description()
}

func (r Response) HttpStatusCode() int {
	return httpStatusCodes[r]
}

func (r Response) OtelCode() codes.Code {
	return otelCodes[r]
}

func AppendDescription(m map[Response]string) {
	for k, v := range m {
		descriptions[k] = v
	}
}

func AppendHttpStatusCode(m map[Response]int) {
	for k, v := range m {
		httpStatusCodes[k] = v
	}
}

func AppendOtelCode(m map[Response]codes.Code) {
	for k, v := range m {
		otelCodes[k] = v
	}
}
