package loginit

import "github.com/sirupsen/logrus"

type MockLogger struct {
	Level   logrus.Level
	Message []interface{}
}

func (l *MockLogger) Print(i ...interface{}) {
	// Not need it yet
}

func (l *MockLogger) Printf(s string, i ...interface{}) {
	// Not need it yet
}

func (l *MockLogger) Println(i ...interface{}) {
	// Not need it yet
}

func (l *MockLogger) Fatal(i ...interface{}) {
	l.Level = logrus.FatalLevel
	l.Message = i
}

func (l *MockLogger) Fatalf(s string, i ...interface{}) {
	l.Level = logrus.FatalLevel
	l.Message = i
}

func (l *MockLogger) Fatalln(i ...interface{}) {
	l.Level = logrus.FatalLevel
	l.Message = i
}

func (l *MockLogger) Panic(i ...interface{}) {
	l.Level = logrus.PanicLevel
	l.Message = i
}

func (l *MockLogger) Panicf(s string, i ...interface{}) {
	l.Level = logrus.PanicLevel
	l.Message = i
}

func (l *MockLogger) Panicln(i ...interface{}) {
	l.Level = logrus.PanicLevel
	l.Message = i
}
