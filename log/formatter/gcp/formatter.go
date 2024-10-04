package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

var (
	logLevel = map[logrus.Level]LogSeverity{
		// PanicLevel level, highest level of severity. Logs and then calls panic with the
		// message passed to Debug, Info, ...
		logrus.PanicLevel: Emergency,
		// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
		// logging level is set to Panic.
		logrus.FatalLevel: Critical,
		// ErrorLevel level. Logs. Used for errors that should definitely be noted.
		// Commonly used for hooks to send errors to an error tracking service.
		logrus.ErrorLevel: Error,
		// WarnLevel level. Non-critical entries that deserve eyes.
		logrus.WarnLevel: Warning,
		// InfoLevel level. General operational entries about what's going on inside the
		// application.
		logrus.InfoLevel: Info,
		// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
		logrus.DebugLevel: Debug,
		// TraceLevel level. Designates finer-grained informational events than the Debug.
		logrus.TraceLevel: Debug,
	}
)

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    },
	// }
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool

	ProjectId string
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+4)
	sourceLocation := make(logrus.Fields)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			switch k {
			case constant.TraceIdLogKey:
				projectId := str.Nil
				if f.ProjectId != str.Empty {
					projectId = f.ProjectId
				}
				data[FieldLoggingGoogleapisComTrace] = fmt.Sprintf(TraceFormat, projectId, v)
			case constant.SpanIdLogKey:
				data[FieldLoggingGoogleapisComSpanID] = v
			case constant.CallerFileLogKey:
				sourceLocation[FieldLoggingGoogleapisComSourceLocationFile] = v
			case constant.CallerFuncLogKey:
				sourceLocation[FieldLoggingGoogleapisComSourceLocationFunction] = v
			case constant.CallerLineLogKey:
				sourceLocation[FieldLoggingGoogleapisComSourceLocationLine] = v
			default:
				data[k] = v
			}
		}
	}

	if f.DataKey != "" {
		newData := make(logrus.Fields, 4)
		newData[f.DataKey] = data
		data = newData
	}

	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339Nano
	}

	if !f.DisableTimestamp {
		data[FieldTimestamp] = entry.Time.Format(timestampFormat)
	}

	data[FieldSeverity] = logLevel[entry.Level]

	if entry.Message != str.Empty {
		data[FieldMessage] = entry.Message
	}

	if len(sourceLocation) > integer.Zero {
		data[FieldLoggingGoogleapisComSourceLocation] = sourceLocation
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!f.DisableHTMLEscape)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return b.Bytes(), nil
}

// This is to not silently overwrite `time`, `msg`, `func` and `level` fields when
// dumping it. If this code wasn't there doing:
//
//	logrus.WithField("level", 1).Info("hello")
//
// Would just silently drop the user provided level. Instead with this code
// it'll logged as:
//
//	{"level": "info", "fields.level": 1, "msg": "hello", "time": "..."}
//
// It's not exported because it's still using Data in an opinionated way. It's to
// avoid code duplication between the two default formatters.
func prefixFieldClashes(data logrus.Fields, fieldMap FieldMap, reportCaller bool) {
	timeKey := fieldMap.resolve(logrus.FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := fieldMap.resolve(logrus.FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := fieldMap.resolve(logrus.FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}

	logrusErrKey := fieldMap.resolve(logrus.FieldKeyLogrusError)
	if l, ok := data[logrusErrKey]; ok {
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	}

	// If reportCaller is not set, 'func' will not conflict.
	if reportCaller {
		funcKey := fieldMap.resolve(logrus.FieldKeyFunc)
		if l, ok := data[funcKey]; ok {
			data["fields."+funcKey] = l
		}
		fileKey := fieldMap.resolve(logrus.FieldKeyFile)
		if l, ok := data[fileKey]; ok {
			data["fields."+fileKey] = l
		}
	}
}
