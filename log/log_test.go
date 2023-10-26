package log

import (
	"context"
	"github.com/rosaekapratama/go-starter/config"
	mocksConfig "github.com/rosaekapratama/go-starter/mocks/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	ctx        context.Context
	oriLog     Logger
	mockLog    *MockLogger
	mockConfig *mocksConfig.MockConfig
)

type LogTestSuite struct {
	suite.Suite
}

func (s *LogTestSuite) SetupSuite() {
	oriLog = logger
}

func (s *LogTestSuite) TearDownSuite() {
	logger = oriLog
}

func (s *LogTestSuite) SetupTest() {
	// Replace logger with mock
	mockLog = &MockLogger{}
	logger = mockLog

	// Init all object
	ctx = context.Background()
	mockConfig = &mocksConfig.MockConfig{}
}

func (s *LogTestSuite) TearDownTest() {
}

func TestLogTestSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}

func (s *LogTestSuite) TestInitInvalidLogLevelError() {
	mockObject := config.Object{Log: &config.LogConfig{Level: "vertigo"}}
	mockConfig.On("GetObject").Return(&mockObject)
	mockLog.EXPECT().Fatalf(mock.Anything, mock.Anything, errInvalidLogLevel, mockObject.Log.Level)
	Init(ctx, mockConfig, "test")
}
