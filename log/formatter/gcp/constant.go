package gcp

const (
	FieldTimestamp                                  string = "timestamp"
	FieldSeverity                                   string = "severity"
	FieldMessage                                    string = "message"
	FieldLoggingGoogleapisComLabels                 string = "logging.googleapis.com/labels"
	FieldLoggingGoogleapisComSourceLocation         string = "logging.googleapis.com/sourceLocation"
	FieldLoggingGoogleapisComSourceLocationFile     string = "file"
	FieldLoggingGoogleapisComSourceLocationLine     string = "line"
	FieldLoggingGoogleapisComSourceLocationFunction string = "function"
	FieldLoggingGoogleapisComInsertID               string = "logging.googleapis.com/insertId"
	FieldLoggingGoogleapisComSpanID                 string = "logging.googleapis.com/spanId"
	FieldLoggingGoogleapisComTrace                  string = "logging.googleapis.com/trace"
	FieldException                                  string = "_exception"
	FieldExceptionFile                              string = "file"
	FieldExceptionMessage                           string = "message"
	FieldExceptionStackTrace                        string = "stackTrace"
)

const (
	Default   LogSeverity = "DEFAULT"
	Debug     LogSeverity = "DEBUG"
	Info      LogSeverity = "INFO"
	Notice    LogSeverity = "NOTICE"
	Warning   LogSeverity = "WARNING"
	Error     LogSeverity = "ERROR"
	Critical  LogSeverity = "CRITICAL"
	Alert     LogSeverity = "ALERT"
	Emergency LogSeverity = "EMERGENCY"
)

const (
	TraceFormat    = "projects/%s/traces/%s"
	DefaultSpanID  = "0000000000000000"
	DefaultTraceID = "00000000000000000000000000000000"
)
