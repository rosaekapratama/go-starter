package config

import (
	"github.com/orandin/lumberjackrus"
	"github.com/rosaekapratama/go-starter/yaml"
)

type Config interface {
	GetObject() *Object
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetBool(key string) (bool, error)
	GetSlice(key string) ([]interface{}, error)
	GetStringAndThrowFatalIfEmpty(key string) string
}

type ConfigImpl struct {
	// All config key and value will be unmarshal here
	o *Object

	// All config key and value in map format will be store here
	m map[string]interface{}
}

type Object struct {
	App           *AppConfig           `yaml:"app"`
	Transport     *TransportConfig     `yaml:"transport"`
	Cors          *CorsConfig          `yaml:"cors"`
	Database      []*DatabaseConfig    `yaml:"database"`
	Redis         *RedisConfig         `yaml:"redis"`
	Log           *LogConfig           `yaml:"log"`
	Otel          *OtelConfig          `yaml:"otel"`
	Google        *GoogleConfig        `yaml:"google"`
	Zeebe         *ZeebeConfig         `yaml:"zeebe"`
	ElasticSearch *ElasticSearchConfig `yaml:"elasticSearch"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

type TransportConfig struct {
	Client *ClientConfig `yaml:"client"`
	Server *ServerConfig `yaml:"server"`
}

type ClientConfig struct {
	Rest *RestClientConfig `yaml:"rest"`
}

type RestClientConfig struct {
	Timeout            int  `yaml:"timeout"`
	InsecureSkipVerify bool `yaml:"insecureSkipVerify"`
}

type ServerConfig struct {
	Rest    *RestServerConfig    `yaml:"rest"`
	Grpc    *GrpcServerConfig    `yaml:"grpc"`
	GraphQL *GraphQLServerConfig `yaml:"graphQL"`
}

type RestServerConfig struct {
	Port *HttpHttpsPortConfig `yaml:"port"`
}

type GrpcServerConfig struct {
	Port *HttpHttpsPortConfig `yaml:"port"`
}

type GraphQLServerConfig struct {
	Port *HttpHttpsPortConfig `yaml:"port"`
}

type HttpHttpsPortConfig struct {
	Http  int `yaml:"http"`
	Https int `yaml:"https"`
}

type CorsConfig struct {
	Pattern          string   `yaml:"pattern"`
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
	Enabled          bool     `yaml:"enabled"`
}

type DatabaseConfig struct {
	Id                        string              `yaml:"id"`
	Driver                    string              `yaml:"driver"`
	Address                   string              `yaml:"address"`
	Database                  string              `yaml:"database"`
	Username                  string              `yaml:"username"`
	Password                  string              `yaml:"password"`
	Conn                      *DatabaseConnConfig `yaml:"conn"`
	SkipDefaultTransaction    bool                `yaml:"skipDefaultTransaction"`
	SlowThreshold             int                 `yaml:"slowThreshold"`
	IgnoreRecordNotFoundError bool                `yaml:"ignoreRecordNotFoundError"`
}

type DatabaseConnConfig struct {
	MaxIdle     int   `yaml:"maxIdle"`
	MaxOpen     int   `yaml:"maxOpen"`
	MaxLifeTime int64 `yaml:"maxLifeTime"`
}

type RedisConfig struct {
	Mode                string `yaml:"mode"`
	RedisSingleConfig   `yaml:",inline"`
	RedisSentinelConfig `yaml:",inline"`
	Disabled            bool `yaml:"disabled"`
}

type RedisSingleConfig struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network *string `yaml:"network"`
	// host:port address.
	Addr string `yaml:"addr"`
	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username *string `yaml:"username"`
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password *string `yaml:"password"`
	// Database to be selected after connecting to the server.
	DB *int `yaml:"db"`
	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries *int `yaml:"maxRetries"`
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff *yaml.Duration `yaml:"minRetryBackoff"`
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff *yaml.Duration `yaml:"maxRetryBackoff"`
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout *yaml.Duration `yaml:"dialTimeout"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout *yaml.Duration `yaml:"readTimeout"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout *yaml.Duration `yaml:"writeTimeout"`

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that fifo has higher overhead compared to lifo.
	PoolFIFO *bool `yaml:"poolFIFO"`
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize *int `yaml:"poolSize"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout *yaml.Duration `yaml:"poolTimeout"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns *int `yaml:"minIdleConns"`
	// Maximum number of idle connections.
	MaxIdleConns *int `yaml:"maxIdleConns"`
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime *yaml.Duration `yaml:"connMaxIdleTime"`
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	ConnMaxLifetime *yaml.Duration `yaml:"connMaxLifetime"`
}

type RedisSentinelConfig struct {
	// The master name.
	MasterName string `yaml:"masterName"`
	// A seed list of host:port addresses of sentinel nodes.
	SentinelAddrs []string `yaml:"sentinelAddrs"`

	// If specified with SentinelPassword, enables ACL-based authentication (via
	// AUTH <user> <pass>).
	SentinelUsername *string `yaml:"sentinelUsername"`
	// Sentinel password from "requirepass <password>" (if enabled) in Sentinel
	// configuration, or, if SentinelUsername is also supplied, used for ACL-based
	// authentication.
	SentinelPassword *string `yaml:"sentinelPassword"`

	// Allows routing read-only commands to the closest master or replica node.
	// This option only works with NewFailoverClusterClient.
	RouteByLatency *bool `yaml:"routeByLatency"`
	// Allows routing read-only commands to the random master or replica node.
	// This option only works with NewFailoverClusterClient.
	RouteRandomly *bool `yaml:"routeRandomly"`

	// Route all commands to replica read-only nodes.
	ReplicaOnly *bool `yaml:"replicaOnly"`

	// Use replicas disconnected with master when cannot get connected replicas
	// Now, this option only works in RandomReplicaAddr function.
	UseDisconnectedReplicas *bool `yaml:"useDisconnectedReplicas"`

	// Following options are copied from Options struct.
	RedisSingleConfig
}

type LogConfig struct {
	Level    string         `yaml:"level"`
	File     *LogFileConfig `yaml:"file"`
	LocalRun bool           `yaml:"localRun"`
}

type LogFileConfig struct {
	lumberjackrus.LogFile
	Enabled bool `yaml:"enabled"`
}

type OtelConfig struct {
	Trace    *OtelTraceConfig  `yaml:"trace"`
	Metric   *OtelMetricConfig `yaml:"metric"`
	Disabled bool              `yaml:"disabled"`
}

type OtelTraceConfig struct {
	Exporter *OtelTraceExporterConfig `yaml:"exporter"`
}

type OtelMetricConfig struct {
	InstrumentationName string                    `yaml:"instrumentationName"`
	Exporter            *OtelMetricExporterConfig `yaml:"exporter"`
}

type OtelTraceExporterConfig struct {
	Type     string                  `yaml:"type"`
	Otlp     *OtelExporterOtlpConfig `yaml:"otlp"`
	Disabled bool                    `yaml:"disabled"`
}

type OtelMetricExporterConfig struct {
	Type     string                  `yaml:"type"`
	Otlp     *OtelExporterOtlpConfig `yaml:"otlp"`
	Disabled bool                    `yaml:"disabled"`
}

type OtelExporterOtlpConfig struct {
	Grpc *OtelExporterOtlpGrpcConfig `yaml:"grpc"`
	Http *OtelExporterOtlpHttpConfig `yaml:"http"`
}

type OtelExporterOtlpGrpcConfig struct {
	Address                     string `yaml:"address"`
	Timeout                     int    `yaml:"timeout"`
	ClientMaxReceiveMessageSize string `yaml:"clientMaxReceiveMessageSize"`
}

type OtelExporterOtlpHttpConfig struct {
	// for future use
}

type GoogleConfig struct {
	Credential string             `yaml:"credential"`
	Cloud      *GoogleCloudConfig `yaml:"cloud"`
	Firebase   *FirebaseConfig    `yaml:"firebase"`
	Disabled   bool               `yaml:"disabled"`
}

type FirebaseConfig struct {
	messaging FirebaseMessagingConfig `yaml:"messaging"`
}

type FirebaseMessagingConfig struct {
	Disabled bool `yaml:"disabled"`
}

type GoogleCloudConfig struct {
	Oauth2 *GoogleCloudOauth2Config `yaml:"oauth2"`
}

type GoogleCloudOauth2Config struct {
	Verification []*GoogleCloudOauth2VerificationConfig `yaml:"verification"`
}

type GoogleCloudOauth2VerificationConfig struct {
	Aud   string `yaml:"aud"`
	Email string `yaml:"email"`
	Sub   string `yaml:"sub"`
}

type ZeebeConfig struct {
	Address                string `yaml:"address"`
	ClientId               string `yaml:"clientId"`
	ClientSecret           string `yaml:"clientSecret"`
	AuthorizationServerURL string `yaml:"authorizationServerURL"`
	Disabled               bool   `yaml:"disabled"`
}

type ElasticSearchConfig struct {
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	Disabled  bool     `yaml:"disabled"`
}
