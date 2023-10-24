package elasticsearch

import (
	"context"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	mocksConfig "github.com/rosaekapratama/go-starter/mocks/config"
	mocksLog "github.com/rosaekapratama/go-starter/mocks/log"
	"github.com/rosaekapratama/go-starter/response"
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

type ElasticSearchTestSuite struct {
	suite.Suite
}

func (s *ElasticSearchTestSuite) SetupSuite() {
	oriLog = log.GetLogger()
}

func (s *ElasticSearchTestSuite) TearDownSuite() {
	log.SetLogger(oriLog)
}

func (s *ElasticSearchTestSuite) SetupTest() {
	// Replace logger with mock
	mockLog = &mocksLog.MockLogger{}
	log.SetLogger(mockLog)

	// Init all object
	ctx = context.Background()
	mockConfig = &mocksConfig.MockConfig{}
}

func (s *ElasticSearchTestSuite) TearDownTest() {
}

func TestElasticSearchTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticSearchTestSuite))
}

func (s *ElasticSearchTestSuite) TestInitCfgAddressEmptyError() {
	mockConfig.On("GetObject").Return(&config.Object{ElasticSearch: &config.ElasticSearchConfig{Disabled: false}})
	mockLog.EXPECT().Fatal(mock.Anything, response.ConfigNotFound, errMissingElasticSearchAddressConfig)
	Init(ctx, mockConfig)
}
