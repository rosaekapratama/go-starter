package log

import (
	"context"
	"github.com/orandin/lumberjackrus"
	"github.com/rosaekapratama/go-starter/constant/env"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/log/formatter/gcp"
	"github.com/rosaekapratama/go-starter/loginit"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/rosaekapratama/go-starter/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

const (
	errInvalidLogLevel = "invalid log level '%s'"
)

var (
	logger Logger
	level  logrus.Level
)

func init() {
	// loginit.Logger assignment need to be put in init(),
	// so this logger can be mocked later in test unit
	logger = &LoggerImpl{logger: loginit.Logger}
}

// Init Set application log
func Init(ctx context.Context, config config.Config, projectId string) {
	cfg := config.GetObject().Log

	// Set logrus configuration
	var err error
	standardLogger := logrus.StandardLogger()
	level, err = logrus.ParseLevel(cfg.Level)
	if err != nil {
		logger.Fatalf(ctx, err, errInvalidLogLevel, cfg.Level)
		return
	}
	var jsonFormatter = gcp.JSONFormatter{ProjectId: projectId}
	if isRunLocally(standardLogger) {
		standardLogger.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			ForceQuote:                true,
			EnvironmentOverrideColors: true,
			FullTimestamp:             true,
		})
	} else {
		standardLogger.SetFormatter(&jsonFormatter)
	}
	standardLogger.SetLevel(level)

	if cfg.File.Enabled {
		// Mkdir log folder
		filePath := cfg.File.Filename
		dir := cfg.GetParentPath()
		if dir != str.Empty {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				logger.Fatalf(ctx, err, "make log dir '%s' failed", dir)
				return
			}
		}

		// Create hook to a file
		hook, err := lumberjackrus.NewHook(
			&lumberjackrus.LogFile{
				Filename:   filePath,
				MaxSize:    cfg.File.MaxSize,
				MaxBackups: cfg.File.MaxBackups,
				MaxAge:     cfg.File.MaxAge,
				Compress:   cfg.File.Compress,
				LocalTime:  cfg.File.LocalTime,
			},
			level,
			&jsonFormatter,
			&lumberjackrus.LogFileOpts{},
		)
		if err != nil {
			logger.Fatal(ctx, err, "log hook creation failed")
			return
		}
		standardLogger.AddHook(hook)

		// Info log file path
		logDir := strings.ReplaceAll(cfg.GetParentPath(), "\\\\", sym.BackSlash)
		logger.Printf(ctx, "logs are printed to '%s'", logDir)
	}

	// Replace logger with configured logger
	logger = &LoggerImpl{logger: standardLogger}
}

func GetLogger() Logger {
	return logger
}

func SetLogger(newLogger Logger) {
	logger = newLogger
}

func addTraceEntries(ctx context.Context, logger logrus.Ext1FieldLogger) logrus.Ext1FieldLogger {
	sc := trace.SpanContextFromContext(ctx)
	newLogger := logger.
		WithField(constant.TraceIdKey, sc.TraceID().String()).
		WithField(constant.SpanIdKey, sc.SpanID().String()).
		WithField(constant.SpanParentIdKey, ctx.Value(constant.SpanParentIdKey))
	return newLogger
}

func addCallerEntries(logger logrus.Ext1FieldLogger) logrus.Ext1FieldLogger {
	if pc, file, line, ok := runtime.Caller(4); ok {
		newLogger := logger.
			WithField(constant.CallerFileKey, file).
			WithField(constant.CallerFuncKey, runtime.FuncForPC(pc).Name()).
			WithField(constant.CallerLineKey, line)
		return newLogger
	}
	return logger
}

// StdEntries Return entries with trace ID entry from span context,
// span ID entry from span context, and
// span parent ID entry from context
func stdEntries(ctx context.Context, logger logrus.Ext1FieldLogger) logrus.Ext1FieldLogger {
	logger = addTraceEntries(ctx, logger)
	logger = addCallerEntries(logger)
	return logger
}

func (logger *LoggerImpl) Trace(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Trace(args...)
}

func (logger *LoggerImpl) Tracef(ctx context.Context, format string, args ...interface{}) {
	stdEntries(ctx, logger.logger).Tracef(format, args...)
}

func (logger *LoggerImpl) Traceln(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Traceln(args...)
}

func (logger *LoggerImpl) Debug(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Debug(args...)
}

func (logger *LoggerImpl) Debugf(ctx context.Context, format string, args ...interface{}) {
	stdEntries(ctx, logger.logger).Debugf(format, args...)
}

func (logger *LoggerImpl) Debugln(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Debugln(args...)
}

func (logger *LoggerImpl) Print(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Print(args...)
}

func (logger *LoggerImpl) Printf(ctx context.Context, format string, args ...interface{}) {
	stdEntries(ctx, logger.logger).Printf(format, args...)
}

func (logger *LoggerImpl) Println(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Println(args...)
}

func (logger *LoggerImpl) Info(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Info(args...)
}

func (logger *LoggerImpl) Infof(ctx context.Context, format string, args ...interface{}) {
	stdEntries(ctx, logger.logger).Infof(format, args...)
}

func (logger *LoggerImpl) Infoln(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Infoln(args...)
}

func (logger *LoggerImpl) Warn(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Warn(args...)
}

func (logger *LoggerImpl) Warnf(ctx context.Context, format string, args ...interface{}) {
	stdEntries(ctx, logger.logger).Warnf(format, args...)
}

func (logger *LoggerImpl) Warnln(ctx context.Context, args ...interface{}) {
	stdEntries(ctx, logger.logger).Warnln(args...)
}

func (logger *LoggerImpl) Error(ctx context.Context, err error, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Error(args...)
}

func (logger *LoggerImpl) Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Errorf(format, args...)
}

func (logger *LoggerImpl) Errorln(ctx context.Context, err error, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Errorln(args...)
}

func (logger *LoggerImpl) Fatal(ctx context.Context, err error, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Fatal(args...)
}

func (logger *LoggerImpl) Fatalf(ctx context.Context, err error, format string, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Fatalf(format, args...)
}

func (logger *LoggerImpl) Fatalln(ctx context.Context, err error, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Fatalln(args...)
}

func (logger *LoggerImpl) Panic(ctx context.Context, err error, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Panic(args...)
}

func (logger *LoggerImpl) Panicf(ctx context.Context, err error, format string, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Panicf(format, args...)
}

func (logger *LoggerImpl) Panicln(ctx context.Context, err error, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
	stdEntries(ctx, logger.logger).WithError(err).Panicln(args...)
}

func (logger *LoggerImpl) GetLevel() logrus.Level {
	return level
}

func (logger *LoggerImpl) GetLogrusLogger() logrus.Ext1FieldLogger {
	return logger.logger
}

func (logger *LoggerImpl) WithField(key string, value interface{}) Logger {
	return &LoggerImpl{logger: logger.logger.WithField(key, value)}
}

func (logger *LoggerImpl) WithFields(fields map[string]interface{}) Logger {
	return &LoggerImpl{logger: logger.logger.WithFields(fields)}
}

func (logger *LoggerImpl) WithTraceFields(ctx context.Context) Logger {
	return &LoggerImpl{logger: addTraceEntries(ctx, logger.logger)}
}

func Trace(ctx context.Context, args ...interface{}) {
	logger.Trace(ctx, args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	logger.Tracef(ctx, format, args...)
}

func Traceln(ctx context.Context, args ...interface{}) {
	logger.Traceln(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	logger.Debug(ctx, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.Debugf(ctx, format, args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	logger.Debugln(ctx, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	logger.Print(ctx, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(ctx, format, args...)
}

func Println(ctx context.Context, args ...interface{}) {
	logger.Println(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	logger.Info(ctx, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.Infof(ctx, format, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	logger.Infoln(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	logger.Warn(ctx, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.Warnf(ctx, format, args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	logger.Warnln(ctx, args...)
}

func Error(ctx context.Context, err error, args ...interface{}) {
	logger.Error(ctx, err, args...)
}

func Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	logger.Errorf(ctx, err, format, args...)
}

func Errorln(ctx context.Context, err error, args ...interface{}) {
	logger.Errorln(ctx, err, args...)
}

func Fatal(ctx context.Context, err error, args ...interface{}) {
	logger.Fatal(ctx, err, args...)
}

func Fatalf(ctx context.Context, err error, format string, args ...interface{}) {
	logger.Fatalf(ctx, err, format, args...)
}

func Fatalln(ctx context.Context, err error, args ...interface{}) {
	logger.Fatalln(ctx, err, args...)
}

func Panic(ctx context.Context, err error, args ...interface{}) {
	logger.Panic(ctx, err, args...)
}

func Panicf(ctx context.Context, err error, format string, args ...interface{}) {
	logger.Panicf(ctx, err, format, args...)
}

func Panicln(ctx context.Context, err error, args ...interface{}) {
	logger.Panicln(ctx, err, args...)
}

func GetLevel() logrus.Level {
	return logger.GetLevel()
}

func GetLogrusLogger() logrus.Ext1FieldLogger {
	return logger.GetLogrusLogger()
}

func WithField(key string, value interface{}) Logger {
	return logger.WithField(key, value)
}

func WithFields(fields map[string]interface{}) Logger {
	return logger.WithFields(fields)
}

func WithTraceFields(ctx context.Context) Logger {
	return logger.WithTraceFields(ctx)
}

func isRunLocally(logger logrus.FieldLogger) bool {
	if localRunStr, ok := os.LookupEnv(env.EnvLocalRun); localRunStr != str.Empty && ok {
		localRun, err := strconv.ParseBool(localRunStr)
		if err != nil {
			logger.Warnf("Failed to parse %s env var '%s' to boolean, %s", env.EnvLocalRun, localRunStr, err.Error())
		} else {
			return localRun
		}
	}

	return false
}
