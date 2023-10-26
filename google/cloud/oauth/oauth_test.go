package oauth

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/log"
	mocksConfig "github.com/rosaekapratama/go-starter/mocks/config"
	mocksLog "github.com/rosaekapratama/go-starter/mocks/log"
	mocksRestClient "github.com/rosaekapratama/go-starter/mocks/transport/restclient"
	"github.com/rosaekapratama/go-starter/transport/restclient"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	ctx         context.Context
	oriLog      log.Logger
	mockLog     *mocksLog.MockLogger
	mockConfig  *mocksConfig.MockConfig
	oriManager  restclient.IManager
	mockManager *mocksRestClient.MockIManager
	client      *restclient.Client
)

type OAuthTestSuite struct {
	suite.Suite
}

func (s *OAuthTestSuite) SetupSuite() {
	oriLog = log.GetLogger()
	oriManager = restclient.Manager
}

func (s *OAuthTestSuite) TearDownSuite() {
	log.SetLogger(oriLog)
	restclient.Manager = oriManager
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

	// Replace manager with mock
	mockManager = &mocksRestClient.MockIManager{}
	restclient.Manager = mockManager

	// Replace default client with mock
	client = &restclient.Client{Resty: resty.New()}
	mockManager.On("GetDefaultClient").Return(client)

	// block all HTTP requests
	httpmock.ActivateNonDefault(client.Resty.GetClient())
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
