package otel

import (
	"context"
	config "github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	mocksConfig "github.com/rosaekapratama/go-starter/mocks/config"
	mocksLog "github.com/rosaekapratama/go-starter/mocks/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	ctx        context.Context
	oriLog     log.Logger
	mockLog    *mocksLog.MockLogger
	mockConfig *mocksConfig.MockConfig
)

type OtelTestSuite struct {
	suite.Suite
}

func (s *OtelTestSuite) SetupSuite() {
	oriLog = log.GetLogger()
}

func (s *OtelTestSuite) TearDownSuite() {
	log.SetLogger(oriLog)
}

func (s *OtelTestSuite) SetupTest() {
	// Replace logger with mock
	mockLog = &mocksLog.MockLogger{}
	log.SetLogger(mockLog)

	// Init all object
	ctx = context.Background()
	mockConfig = &mocksConfig.MockConfig{}
}

func (s *OtelTestSuite) TearDownTest() {
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(OtelTestSuite))
}

func (s *OtelTestSuite) TestInitIsDisabled() {
	mockObject := config.Object{Otel: &config.OtelConfig{
		Disabled: true,
	}}
	mockConfig.On("GetObject").Return(&mockObject)
	mockLog.EXPECT().Info(mock.Anything, errOtelInitIsDisabled)
	Init(ctx, mockConfig)
}
