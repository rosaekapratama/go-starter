package avro

import (
	"context"
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

type AvroTestSuite struct {
	suite.Suite
}

func (s *AvroTestSuite) SetupSuite() {
	oriLog = log.GetLogger()
}

func (s *AvroTestSuite) TearDownSuite() {
	// Recover original logger
	log.SetLogger(oriLog)
}

func (s *AvroTestSuite) SetupTest() {
	// Replace logger with mock
	mockLog = &mocksLog.MockLogger{}
	log.SetLogger(mockLog)

	// Init all object
	ctx = context.Background()
	mockConfig = &mocksConfig.MockConfig{}
}

func (s *AvroTestSuite) TearDownTest() {
}

func TestAvroTestSuite(t *testing.T) {
	suite.Run(t, new(AvroTestSuite))
}

func (s *AvroTestSuite) TestInitReadAvroDirPathError() {
	mockLog.EXPECT().Warn(mock.Anything, errAvroSchemaManagerIsDisabled)
	Init(ctx)
}
