package response

import (
	"fmt"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
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
	DBConnIdNotFound
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
	FileSizeMustBeGreaterThanZero
	PageSizeExceedsMaxLimit
)

var (
	descriptionMap = map[Response]string{
		Success:                       "success",
		GeneralError:                  "general error",
		InitFailed:                    "init failed",
		InvalidConfigValueType:        "invalid config value type",
		InvalidArgument:               "invalid argument",
		InvalidPageRequest:            "invalid page request",
		DataNotFound:                  "data not found",
		ConfigNotFound:                "config not found",
		DBConnIdNotFound:              "database connection ID not found",
		UnknownResponse:               "unknown response",
		APINotRegistered:              "api is not registered, path=%s, method=%s",
		DataIsEmpty:                   "data is empty",
		UnauthorizedAccess:            "unauthorized access",
		InvalidBodyRequest:            "invalid body request",
		OperationNotPermitted:         "operation not permitted",
		RequestEntityTooLarge:         "request entity too large",
		InvalidCredentials:            "invalid credentials",
		AvroSchemaNotFound:            "avro schema not found",
		InvalidConfig:                 "invalid config",
		Processing:                    "processing",
		UnsupportedType:               "unsupported type",
		MissingGoogleToken:            "missing google token",
		FileSizeMustBeGreaterThanZero: "file size must be greater than zero",
		PageSizeExceedsMaxLimit:       "page size exceeds the maximum limit of 100",
	}

	httpStatusCodeMap = map[Response]int{
		Success:                       http.StatusOK,
		GeneralError:                  http.StatusInternalServerError,
		InitFailed:                    http.StatusInternalServerError,
		InvalidConfigValueType:        http.StatusInternalServerError,
		InvalidArgument:               http.StatusBadRequest,
		InvalidPageRequest:            http.StatusBadRequest,
		DataNotFound:                  http.StatusOK,
		ConfigNotFound:                http.StatusInternalServerError,
		DBConnIdNotFound:              http.StatusInternalServerError,
		UnknownResponse:               http.StatusInternalServerError,
		APINotRegistered:              http.StatusNotFound,
		DataIsEmpty:                   http.StatusOK,
		UnauthorizedAccess:            http.StatusUnauthorized,
		InvalidBodyRequest:            http.StatusBadRequest,
		OperationNotPermitted:         http.StatusForbidden,
		RequestEntityTooLarge:         http.StatusRequestEntityTooLarge,
		InvalidCredentials:            http.StatusUnauthorized,
		AvroSchemaNotFound:            http.StatusInternalServerError,
		InvalidConfig:                 http.StatusInternalServerError,
		Processing:                    http.StatusProcessing,
		UnsupportedType:               http.StatusOK,
		MissingGoogleToken:            http.StatusBadRequest,
		FileSizeMustBeGreaterThanZero: http.StatusBadRequest,
		PageSizeExceedsMaxLimit:       http.StatusBadRequest,
	}

	otelCodeMap = map[Response]otelCodes.Code{
		Success:                       otelCodes.Ok,
		GeneralError:                  otelCodes.Error,
		InitFailed:                    otelCodes.Error,
		InvalidConfigValueType:        otelCodes.Error,
		InvalidArgument:               otelCodes.Ok,
		InvalidPageRequest:            otelCodes.Ok,
		DataNotFound:                  otelCodes.Ok,
		DBConnIdNotFound:              otelCodes.Error,
		ConfigNotFound:                otelCodes.Error,
		UnknownResponse:               otelCodes.Error,
		APINotRegistered:              otelCodes.Ok,
		DataIsEmpty:                   otelCodes.Ok,
		UnauthorizedAccess:            otelCodes.Ok,
		InvalidBodyRequest:            otelCodes.Ok,
		OperationNotPermitted:         otelCodes.Ok,
		RequestEntityTooLarge:         otelCodes.Ok,
		InvalidCredentials:            otelCodes.Ok,
		AvroSchemaNotFound:            otelCodes.Ok,
		InvalidConfig:                 otelCodes.Ok,
		Processing:                    otelCodes.Ok,
		UnsupportedType:               otelCodes.Ok,
		MissingGoogleToken:            otelCodes.Ok,
		FileSizeMustBeGreaterThanZero: otelCodes.Ok,
		PageSizeExceedsMaxLimit:       otelCodes.Ok,
	}

	grpcCodeMap = map[Response]grpcCodes.Code{
		Success:                       grpcCodes.OK,
		GeneralError:                  grpcCodes.Unknown,
		InitFailed:                    grpcCodes.FailedPrecondition,
		InvalidConfigValueType:        grpcCodes.FailedPrecondition,
		InvalidArgument:               grpcCodes.InvalidArgument,
		InvalidPageRequest:            grpcCodes.InvalidArgument,
		DataNotFound:                  grpcCodes.NotFound,
		DBConnIdNotFound:              grpcCodes.Unknown,
		ConfigNotFound:                grpcCodes.Unknown,
		UnknownResponse:               grpcCodes.Unknown,
		APINotRegistered:              grpcCodes.Unimplemented,
		DataIsEmpty:                   grpcCodes.NotFound,
		UnauthorizedAccess:            grpcCodes.PermissionDenied,
		InvalidBodyRequest:            grpcCodes.InvalidArgument,
		OperationNotPermitted:         grpcCodes.PermissionDenied,
		RequestEntityTooLarge:         grpcCodes.InvalidArgument,
		InvalidCredentials:            grpcCodes.PermissionDenied,
		AvroSchemaNotFound:            grpcCodes.FailedPrecondition,
		InvalidConfig:                 grpcCodes.FailedPrecondition,
		Processing:                    grpcCodes.OK,
		UnsupportedType:               grpcCodes.InvalidArgument,
		MissingGoogleToken:            grpcCodes.FailedPrecondition,
		FileSizeMustBeGreaterThanZero: grpcCodes.InvalidArgument,
		PageSizeExceedsMaxLimit:       grpcCodes.InvalidArgument,
	}

	isErrorMap = map[Response]bool{
		Success:                       false,
		GeneralError:                  true,
		InitFailed:                    true,
		InvalidConfigValueType:        true,
		InvalidArgument:               false,
		InvalidPageRequest:            false,
		DataNotFound:                  false,
		DBConnIdNotFound:              true,
		ConfigNotFound:                true,
		UnknownResponse:               true,
		APINotRegistered:              true,
		DataIsEmpty:                   false,
		UnauthorizedAccess:            false,
		InvalidBodyRequest:            true,
		OperationNotPermitted:         false,
		RequestEntityTooLarge:         false,
		InvalidCredentials:            false,
		AvroSchemaNotFound:            true,
		InvalidConfig:                 true,
		Processing:                    false,
		UnsupportedType:               false,
		MissingGoogleToken:            true,
		FileSizeMustBeGreaterThanZero: false,
		PageSizeExceedsMaxLimit:       false,
	}
)

func (r Response) Code() string {
	return fmt.Sprintf("%04d", r)
}

func (r Response) Description() string {
	return descriptionMap[r]
}

func (r Response) Error() string {
	return r.Description()
}

func (r Response) HttpStatusCode() int {
	return httpStatusCodeMap[r]
}

func (r Response) OtelCode() otelCodes.Code {
	return otelCodeMap[r]
}

func (r Response) GrpcCode() grpcCodes.Code {
	if v, exists := grpcCodeMap[r]; exists {
		return v
	} else {
		return 2
	}
}

func (r Response) IsError() bool {
	if v, exists := isErrorMap[r]; exists {
		return v
	} else {
		return true
	}
}

func AppendDescription(m map[Response]string) {
	for k, v := range m {
		descriptionMap[k] = v
	}
}

func AppendHttpStatusCode(m map[Response]int) {
	for k, v := range m {
		httpStatusCodeMap[k] = v
	}
}

func AppendOtelCode(m map[Response]otelCodes.Code) {
	for k, v := range m {
		otelCodeMap[k] = v
	}
}

func AppendGrpcCode(m map[Response]grpcCodes.Code) {
	for k, v := range m {
		grpcCodeMap[k] = v
	}
}

func GrpcErrFromErr(err error) (grpcErr error) {
	if err != nil {
		switch v := err.(type) {
		case Response:
			return grpcStatus.Error(v.GrpcCode(), v.Description())
		default:
			return grpcStatus.Error(grpcCodes.Unknown, err.Error())
		}
	}
	return
}
