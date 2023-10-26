package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Trace(ctx context.Context, args ...interface{})
	Tracef(ctx context.Context, format string, args ...interface{})
	Traceln(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Debugln(ctx context.Context, args ...interface{})
	Print(ctx context.Context, args ...interface{})
	Printf(ctx context.Context, format string, args ...interface{})
	Println(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Infoln(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warnln(ctx context.Context, args ...interface{})
	Error(ctx context.Context, err error, args ...interface{})
	Errorf(ctx context.Context, err error, format string, args ...interface{})
	Errorln(ctx context.Context, err error, args ...interface{})
	Fatal(ctx context.Context, err error, args ...interface{})
	Fatalf(ctx context.Context, err error, format string, args ...interface{})
	Fatalln(ctx context.Context, err error, args ...interface{})
	Panic(ctx context.Context, err error, args ...interface{})
	Panicf(ctx context.Context, err error, format string, args ...interface{})
	Panicln(ctx context.Context, err error, args ...interface{})

	GetLevel() logrus.Level
	GetLogrusLogger() logrus.Ext1FieldLogger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithTraceFields(ctx context.Context) Logger
}

type LoggerImpl struct {
	logger logrus.Ext1FieldLogger
}
