package gcp

type LogSeverity string

type Layout struct {
	Timestamp                          string                                    `json:"timestamp"`
	Severity                           string                                    `json:"severity"`
	Message                            string                                    `json:"message,omitempty"`
	HttpRequest                        *HttpRequest                              `json:"httpRequest,omitempty"`
	LoggingGoogleapisComLabels         string                                    `json:"logging.googleapis.com/labels,omitempty"`
	LoggingGoogleapisComSourceLocation *LoggingGoogleapisComSourceLocationLayout `json:"logging.googleapis.com/sourceLocation,omitempty"`
	LoggingGoogleapisComInsertID       string                                    `json:"logging.googleapis.com/insertId,omitempty"`
	LoggingGoogleapisComSpanID         string                                    `json:"logging.googleapis.com/spanId"`
	LoggingGoogleapisComTrace          string                                    `json:"logging.googleapis.com/trace"`
	Exception                          *ExceptionLayout                          `json:"_exception,omitempty"`
}

type LoggingGoogleapisComSourceLocationLayout struct {
	File     string `json:"file,omitempty"`
	Line     string `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}

type ExceptionLayout struct {
	File       string `json:"file,omitempty"`
	Message    string `json:"message,omitempty"`
	StackTrace string `json:"stackTrace,omitempty"`
}

type HttpRequest struct {
	RequestMethod                  string `json:"requestMethod,omitempty"`
	RequestUrl                     string `json:"requestUrl,omitempty"`
	RequestSize                    string `json:"requestSize,omitempty"`
	Status                         int    `json:"status,omitempty"`
	ResponseSize                   string `json:"responseSize,omitempty"`
	UserAgent                      string `json:"userAgent,omitempty"`
	RemoteIp                       string `json:"remoteIp,omitempty"`
	ServerIp                       string `json:"serverIp,omitempty"`
	Referer                        string `json:"referer,omitempty"`
	Latency                        string `json:"latency,omitempty"`
	CacheLookup                    bool   `json:"cacheLookup,omitempty"`
	CacheHit                       bool   `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 string `json:"cacheFillBytes,omitempty"`
	Protocol                       string `json:"protocol,omitempty"`
}
