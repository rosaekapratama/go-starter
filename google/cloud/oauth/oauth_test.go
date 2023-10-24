package oauth

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	mocksConfig "github.com/rosaekapratama/go-starter/mocks/config"
	mocksLog "github.com/rosaekapratama/go-starter/mocks/log"
	"github.com/rosaekapratama/go-starter/transport/rest/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	ctx           context.Context
	oriLog        log.Logger
	mockLog       *mocksLog.MockLogger
	mockConfig    *mocksConfig.MockConfig
	oriHttpClient = client.GetDefaultClient()
)

type OAuthTestSuite struct {
	suite.Suite
}

func (s *OAuthTestSuite) SetupSuite() {
	oriLog = log.GetLogger()
}

func (s *OAuthTestSuite) TearDownSuite() {
	log.SetLogger(oriLog)
}

func (s *OAuthTestSuite) SetupTest() {
	// Replace logger with mock
	mockLog = &mocksLog.MockLogger{}
	log.SetLogger(mockLog)

	// Init all object
	ctx = context.Background()
	mockConfig = &mocksConfig.MockConfig{}
	mockConfig.On("GetObject").Return(
		&config.Object{
			Transport: &config.TransportConfig{
				Client: &config.ClientConfig{
					Rest: &config.RestClientConfig{
						InsecureSkipVerify: false,
					},
				},
			},
		},
	)
	client.Init(mockConfig)
	httpClient = client.New()
	// block all HTTP requests
	httpmock.ActivateNonDefault(httpClient.Resty.GetClient())
}

func (s *OAuthTestSuite) TearDownTest() {
}

func TestOAuthTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthTestSuite))
}

func (s *OAuthTestSuite) TestVerifyTokenHttpGetError() {
	oauthClient := &ClientImpl{tokeninfoEndpoint: "test"}
	mockLog.EXPECT().Error(mock.Anything, mock.Anything, errFailedToGetTokenInfo)
	_, _ = oauthClient.VerifyToken(ctx, "token")
}
