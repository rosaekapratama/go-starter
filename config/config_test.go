package config

import (
	"github.com/go-playground/assert/v2"
	"github.com/rosaekapratama/go-starter/loginit"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

var (
	oriLog             = logger
	oriArgs            = os.Args
	mockLog            = loginit.MockLogger{}
	mockConfigFileName = "app-mock.yaml"
	yamlConfig         = `
app:
  name: YourAppName
  mode: Production

transport:
  client:
    rest:
      insecureSkipVerify: false

  server:
    rest:
      port:
        http: 8080
        https: 8443
    grpc:
      port:
        http: 50051
        https: 5443
    graphQL:
      port:
        http: 8081
        https: 8444

cors:
  pattern: "*"
  allowOrigins:
    - http://example.com
    - https://example.com
  allowMethods:
    - GET
    - POST
  allowHeaders:
    - Content-Type
  exposeHeaders:
    - X-Custom-Header
  allowCredentials: true
  maxAge: 3600
  disabled: false

database:
  db1:
    driver: mysql
    address: localhost:3306
    database: mydb
    username: user1
    password: secret1
    conn:
      maxIdle: 10
      maxOpen: 100
      maxLifeTime: 3600
  db2:
    driver: postgres
    address: localhost:5432
    database: anotherdb
    username: user2
    password: secret2
    conn:
      maxIdle: 5
      maxOpen: 50
      maxLifeTime: 1800

redis:
  mode: single
  network: tcp
  addr: localhost:6379
  username: redisuser
  password: redispassword
  db: 0
  maxRetries: 3
  minRetryBackoff: 8ms
  maxRetryBackoff: 512ms
  dialTimeout: 5s
  readTimeout: 3s
  writeTimeout: 3s
  poolFIFO: true
  poolSize: 10
  poolTimeout: 4s
  minIdleConns: 5
  maxIdleConns: 10
  connMaxIdleTime: 300s
  connMaxLifetime: 0s
  disabled: false

log:
  level: info
  file:
    LogFile:
      Filename: /var/log/yourapp.log
      MaxSize: 100
      MaxBackups: 3
      MaxAge: 28
      Compress: true
    enabled: true
  localRun: false

otel:
  trace:
    exporter:
      type: otlp
      otlp:
        grpc:
          address: localhost:4317
          timeout: 5
          clientMaxReceiveMessageSize: 1048576
        disabled: false

  metric:
    instrumentationName: your-app
    exporter:
      type: otlp
      otlp:
        grpc:
          address: localhost:4318
          timeout: 5
          clientMaxReceiveMessageSize: 1048576
        disabled: false
  disabled: false

google:
  credential: /path/to/google/credential.json
  cloud:
    oauth2:
      verification:
        - aud: audience1
          email: email1@example.com
          sub: subject1
        - aud: audience2
          email: email2@example.com
          sub: subject2

firebase:
  messaging:
    disabled: false

zeebe:
  address: localhost:26500
  clientId: yourClientId
  clientSecret: yourClientSecret
  authorizationServerURL: https://your-auth-server.com

elasticSearch:
  addresses:
    - http://localhost:9200
  username: esuser
  password: espassword
`
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) SetupSuite() {
	// Create valid mock config file
	err := os.WriteFile(mockConfigFileName, []byte(yamlConfig), 0755)
	if err != nil {
		oriLog.Fatal(err)
	}

	// Provider valid mock config file to os arguments
	os.Args = []string{"", mockConfigFileName}
}

func (s *ConfigTestSuite) TearDownSuite() {
	// Delete valid mock config file
	err := os.Remove(mockConfigFileName)
	if err != nil {
		oriLog.Fatal(err)
	}

	// Recover original os arguments
	os.Args = oriArgs
}

func (s *ConfigTestSuite) SetupTest() {
	// Replace logger with mock
	logger = &mockLog
}

func (s *ConfigTestSuite) TearDownTest() {
	// Recover original logger
	logger = oriLog
	os.Args = oriArgs
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestInitReadDefaultConfigFileError() {
	// Keep current os arguments
	temp := os.Args

	// Replace os arguments with empty one so it will use default path config
	os.Args = make([]string, 0)

	Init()
	assert.Equal(s.T(), mockLog.Level, logrus.FatalLevel)
	assert.Equal(s.T(), mockLog.Message[0], errReadingConfigFile)

	// Recover os args for other unit test
	os.Args = temp
}

func (s *ConfigTestSuite) TestInitUnmarshalConfigFileError() {
	// Create invalid mock config file
	cfgStr := strings.ReplaceAll(yamlConfig, "name", "")
	fn := "temp.yaml"
	err := os.WriteFile(fn, []byte(cfgStr), 0755)
	if err != nil {
		oriLog.Fatal(err)
	}
	os.Args = []string{"", fn}

	Init()
	defer func() {
		err = os.Remove(fn)
		if err != nil {
			oriLog.Fatal(err)
		}
	}()
	assert.Equal(s.T(), mockLog.Level, logrus.FatalLevel)
	assert.Equal(s.T(), mockLog.Message[0], errUnmarshalConfigFile)
}

func (s *ConfigTestSuite) TestAppConfigIsNil() {
	// Create missing app mock config file
	cfgStr := strings.ReplaceAll(yamlConfig, "app:\n  name: YourAppName\n  mode: Production", "")
	fn := "temp.yaml"
	err := os.WriteFile(fn, []byte(cfgStr), 0755)
	if err != nil {
		oriLog.Fatal(err)
	}
	os.Args = []string{"", fn}

	Init()
	defer func() {
		err = os.Remove(fn)
		if err != nil {
			oriLog.Fatal(err)
		}
	}()
	assert.Equal(s.T(), mockLog.Level, logrus.FatalLevel)
	assert.Equal(s.T(), mockLog.Message[0], errMissingApplicationName)
}

func (s *ConfigTestSuite) TestAppNameConfigIsEmpty() {
	// Create missing app mock config file
	cfgStr := strings.ReplaceAll(yamlConfig, "app:\n  name: YourAppName\n  mode: Production", "app:\n  name: ''\n  mode: debug")
	fn := "temp.yaml"
	err := os.WriteFile(fn, []byte(cfgStr), 0755)
	if err != nil {
		oriLog.Fatal(err)
	}
	os.Args = []string{"", fn}

	Init()
	defer func() {
		err = os.Remove(fn)
		if err != nil {
			oriLog.Fatal(err)
		}
	}()
	assert.Equal(s.T(), mockLog.Level, logrus.FatalLevel)
	assert.Equal(s.T(), mockLog.Message[0], errMissingApplicationName)
}
