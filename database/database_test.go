package database

import (
	"context"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	mocksConfig "github.com/rosaekapratama/go-starter/mocks/config"
	mocksLog "github.com/rosaekapratama/go-starter/mocks/log"
	"github.com/sirupsen/logrus"
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

type DatabaseTestSuite struct {
	suite.Suite
}

func (s *DatabaseTestSuite) SetupSuite() {
	oriLog = log.GetLogger()
}

func (s *DatabaseTestSuite) TearDownSuite() {
	// Recover original logger
	log.SetLogger(oriLog)
}

func (s *DatabaseTestSuite) SetupTest() {
	// Replace logger with mock
	mockLog = &mocksLog.MockLogger{}
	log.SetLogger(mockLog)

	// Init all object
	ctx = context.Background()
	mockConfig = &mocksConfig.MockConfig{}
}

func (s *DatabaseTestSuite) TearDownTest() {
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) TestInitIsDisabled() {
	mockConfig.On("GetObject").Return(&config.Object{Database: nil})
	mockLog.On("GetLevel").Return(logrus.WarnLevel)
	mockLog.EXPECT().Warn(mock.Anything, errDatabaseManagerIsDisabled)
	Init(ctx, mockConfig)
}
