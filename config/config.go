package config

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/orandin/lumberjackrus"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/loginit"
	"github.com/rosaekapratama/go-starter/response"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

const (
	errReadingConfigFile      = "error reading config file, "
	errUnmarshalConfigFile    = "error unmarshal config file to instance, "
	errMissingApplicationName = "missing application name"
	errFailedToGetConfig      = "failed to get config, configKey=%s, error=%s"
	errConfigValueIsEmpty     = "config value must not empty, configKey=%s"
)

var (
	logger   logrus.StdLogger
	Instance Config
)

func init() {
	// loginit.Logger assignment need to be put in init(),
	// so this logger can be mocked later in test unit
	logger = loginit.Logger
}

// Init Set config object from file path and store it to singleton
func Init() {
	// Get file path from argument
	// set default if not exists
	var filePath string
	args := os.Args
	if len(args) < 2 {
		filePath = defaultFilePath
	} else {
		filePath = args[1]
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		logger.Fatal(errReadingConfigFile, err)
		return
	}

	// Set default value for all fields
	o := &Object{
		App: &AppConfig{Mode: gin.DebugMode},
		Transport: &TransportConfig{
			Client: &ClientConfig{
				Rest: &RestClientConfig{InsecureSkipVerify: false},
			},
			Server: &ServerConfig{
				Rest: &RestServerConfig{
					Port: &HttpHttpsPortConfig{
						Http:  defaultHTTPRESTPort,
						Https: defaultHTTPSRESTPort,
					},
				},
				Grpc: &GrpcServerConfig{
					Port: &HttpHttpsPortConfig{
						Http:  defaultHTTPGRPCPort,
						Https: defaultHTTPSGRPCPort,
					},
				},
				GraphQL: &GraphQLServerConfig{
					Port: &HttpHttpsPortConfig{
						Http:  defaultHTTPGraphQLPort,
						Https: defaultHTTPSGraphQLPort,
					},
				},
			},
		},
		Log: &LogConfig{
			Level: "info",
			File: &LogFileConfig{
				LogFile: lumberjackrus.LogFile{
					MaxSize: 100,
					MaxAge:  30,
				},
				Enabled: false,
			},
		},
		Otel: &OtelConfig{
			Trace: &OtelTraceConfig{
				Exporter: &OtelTraceExporterConfig{
					Otlp: &OtelExporterOtlpConfig{
						Grpc: &OtelExporterOtlpGrpcConfig{
							Timeout: 30,
						},
					},
					Disabled: true,
				},
			},
			Metric: &OtelMetricConfig{
				Exporter: &OtelMetricExporterConfig{
					Otlp: &OtelExporterOtlpConfig{
						Grpc: &OtelExporterOtlpGrpcConfig{
							Timeout: 30,
						},
					},
					Disabled: true,
				},
			},
			Disabled: true,
		},
		Google: &GoogleConfig{
			Cloud: &GoogleCloudConfig{
				Pubsub: &GoogleCloudPubsubConfig{
					Publisher: &GoogleCloudPubsubPublisherConfig{
						Logging: false,
					},
					Subscriber: &GoogleCloudPubsubSubscriberConfig{
						Logging: false,
					},
				},
			},
		},
	}

	// Expand environment variables
	bytes = []byte(os.ExpandEnv(string(bytes)))

	// Unmarshal and set value from config file
	err = yaml.Unmarshal(bytes, o)
	if err != nil {
		logger.Fatal(errUnmarshalConfigFile, err)
		return
	}

	// Validate mandatory config
	if o.App == nil || o.App.Name == str.Empty {
		logger.Fatal(errMissingApplicationName)
		return
	}

	// Set default file log name if not exists
	if o.Log.File.Filename == str.Empty {
		s := fmt.Sprintf("%c", os.PathSeparator)
		sb := strings.Builder{}
		sb.WriteString(os.TempDir())
		sb.WriteString(s)
		sb.WriteString("app")
		sb.WriteString(s)
		sb.WriteString(o.App.Name)
		sb.WriteString(s)
		sb.WriteString(o.App.Name)
		sb.WriteString(".log")
		o.Log.File.Filename = sb.String()
	}

	// Set default value for open telemetry metric if no config defined
	if o.Otel.Metric == nil {
		sb := strings.Builder{}
		sb.WriteString(o.App.Name)
		sb.WriteString("-meter")
		instrumentationName := sb.String()
		o.Otel.Metric = &OtelMetricConfig{
			InstrumentationName: instrumentationName,
		}
	}

	// Unmarshal to map
	m := structs.Map(o)
	err = yaml.Unmarshal(bytes, m)
	if err != nil {
		logger.Fatal("error unmarshal config file to map, ", err)
		return
	}

	Instance = &ConfigImpl{
		o: o,
		m: m,
	}
}

func (c *ConfigImpl) GetObject() *Object {
	return c.o
}

func getVal(key string, m map[string]interface{}) interface{} {
	if key == str.Empty {
		return nil
	}

	keys := strings.SplitN(key, sym.Dot, integer.Two)
	if v, ok := m[keys[integer.Zero]]; ok {
		switch v := v.(type) {
		case map[string]interface{}:
			return getVal(keys[integer.One], v)
		default:
			return v
		}
	}
	return nil
}

// GetString use dot to get value from nested key
// ex: keycloak.address
func (c *ConfigImpl) GetString(key string) (string, error) {
	v := getVal(key, c.m)
	if v == nil {
		return str.Empty, nil
	}

	switch v := v.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return str.Empty, response.InvalidConfigValueType
	}
}

func (c *ConfigImpl) GetStringAndThrowFatalIfEmpty(key string) string {
	v, err := c.GetString(key)
	if err != nil {
		logger.Fatal(fmt.Sprintf(errFailedToGetConfig, key, err.Error()))
		return str.Empty
	}
	if v == str.Empty {
		logger.Fatal(fmt.Sprintf(errConfigValueIsEmpty, key))
		return str.Empty
	}
	return v
}

// GetInt use dot to get value from nested key, ex: keycloak.address
func (c *ConfigImpl) GetInt(key string) (int, error) {
	v := getVal(key, c.m)
	if v == nil {
		return integer.Zero, response.ConfigNotFound
	}

	switch v := v.(type) {
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return integer.Zero, err
		}
		return i, nil
	case int:
		return v, nil
	case bool:
		if v {
			return integer.One, nil
		}
		return integer.Zero, nil
	default:
		return integer.Zero, response.InvalidConfigValueType
	}
}

// GetBool use dot to get value from nested key, ex: keycloak.address
func (c *ConfigImpl) GetBool(key string) (bool, error) {
	v := getVal(key, c.m)
	if v == nil {
		return false, response.ConfigNotFound
	}

	switch v := v.(type) {
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, err
		}
		return b, nil
	case int:
		i := v
		if i == integer.One {
			return true, nil
		} else if i == integer.Zero {
			return false, nil
		}
		return false, response.InvalidConfigValueType
	case bool:
		return v, nil
	default:
		return false, response.InvalidConfigValueType
	}
}

// GetSlice use dot to get value from nested key, ex: keycloak.address
func (c *ConfigImpl) GetSlice(key string) ([]interface{}, error) {
	v := getVal(key, c.m)
	if v == nil {
		return nil, nil
	}

	switch v := v.(type) {
	case []interface{}:
		return v, nil
	default:
		return nil, response.InvalidConfigValueType
	}
}

func (c *LogConfig) GetParentPath() string {
	filePath := c.File.Filename
	s := fmt.Sprintf("%c", os.PathSeparator)
	i := strings.LastIndex(filePath, s)
	if i < 0 {
		return str.Empty
	}
	return filePath[:i+1]
}
