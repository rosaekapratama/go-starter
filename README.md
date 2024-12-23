# go-starter #

Auto configuration and collection of modules for Go REST services powered by Gin web framework and GORM, now support GRPC also.

## Usage

### Installation

```shell
go get github.com/rosaekapratama/go-starter
```

### Importing

```go
import "github.com/rosaekapratama/go-starter"
```

### Configuration File ###
```yaml
---
app:
  name: my-service
  mode: debug # debug/release/test, for production use release

transport:
  client:
    rest:
      logging:
        payloadLogSizeLimit: 2KB
        stdout: false # Show incoming and outgoing message globally, default is false
        database: pgsql1 # Write log to database with ID pgsql1
      timeout: 30 # Wait time in second
      insecure: false
    soap:
      logging:
        payloadLogSizeLimit: 2KB
        stdout: false # Show incoming and outgoing message globally, default is false
        database: pgsql1 # Write log to database with ID pgsql1
      timeout: 30 # Wait time in second
      insecure: false
    grpc:
      connId1:
        address: localhost:8080
        timeout: 30
        insecure: false
        logging:
          payloadLogSizeLimit: 2KB
          stdout: false # Show incoming and outgoing message globally, default is false
          database: pgsql1 # Write log to database with ID pgsql1
    ssh:
      myClientId:
        address: localhost:22
        username: admin
        auth: # Choose 1 between below option, password/privateKey/privateKeyPath
          password: admin123
          privateKey: "this should be private key in raw or base64"
          privateKeyPath: "/tmp/privateKeyLocation"
        idleTimeout: 5000 # connection will be closed if no operation in millisecond
  server:
    rest:
      logging:
        payloadLogSizeLimit: 2KB
        stdout: false # Show incoming and outgoing message globally, default is false
        database: pgsql1 # Write log to database with ID pgsql1
      port:
        http: 9092
    grpc:
      port:
        http: 9092
      logging:
        payloadLogSizeLimit: 2KB
        stdout: false # Show incoming and outgoing message globally, default is false
        database: pgsql1 # Write log to database with ID pgsql1
      disabled: true # true will disable the GRPC server

cors:
  pattern: "*"
  allowOrigins:
    - "*"
  allowMethods:
    - "GET"
    - "POST"
    - "DELETE"
    - "PATCH"
  allowHeaders:
    - "Origin"
    - "Authorization"
    - "Access-Control-Allow-Origin"
    - "Access-Control-Allow-Headers"
    - "Content-Type"
    - "Page-Num"
    - "Page-Size"
    - "Realm"
  exposeHeaders:
    - "Access-Control-Allow-Origin"
    - "Access-Control-Allow-Headers"
    - "Content-Type"
    - "Content-Length"
    - "Content-Description"
    - "Content-Disposition"
  allowCredentials: true # Credentials share
  maxAge: 43200 # Preflight requests cached for 12 hours

log:
  # Valid level are
  # trace, debug, info, warn, error, fatal, panic.
  level: info

  # Please refer to lumberjackrus.LogFile struct to set log.file.xxx fields.
  # Because of gopkg.in/yaml behaviour,
  # all of its field below must be lowercased.
  file:
    filename: my-service.log # If not set, file name will be set to os.TempDir/app/[app.appname]/[app.appname].log
    maxsize: 100
    maxage: 3
    maxbackups: 3
    localtime: false
    compress: false
    enable: true

# Choose mode between single or sentinel
redis:
  mode: sentinel
  masterName: mymaster
  sentinelAddrs:
    - master:6379
    - slave1:6379
    - slave2:6379
  sentinelPassword: ''
  addr: localhost:6379
  password: ''
  db: 0
  disabled: false # If true then redis init will be skipped
  
  # Other config
  network: 'tcp' # TCP/UNIX, default is TCP
  username: ''
  maxRetries: 3 # -1 (not 0) disables retries. 
  minRetryBackoff: 8 # In milliseconds
  maxRetryBackoff: 512 # In milliseconds
  dialTimeout: 5 # In seconds
  readTimeout: 3 # In seconds
  writeTimeout: 3 # Default is readTimeout in seconds
  poolFIFO: false
  poolSize: 10
  PoolTimeout: 4 # Default is readTimeout + 1 second
  minIdleConns: 10
  maxIdleConns: 10
  connMaxIdleTime: 300 # In seconds, default is 5 minutes. -1 disables idle timeout check
  connMaxLifetime: -1
  
  # Other sentinel config
  routeByLatency: true
  routeRandomly: false
  replicaOnly: false
  useDisconnectedReplicas: false

database:
  pgsql1:
    driver: pgsql
    address: localhost
    database: mydatabase
    username: myuser
    password: mypass
    conn:
      maxIdle: 30
      maxOpen: 30
      maxLifeTime: 60000
    slowThreshold: 200 # In millisecond
    skipDefaultTransaction: false
    ignoreRecordNotFoundError: false
  pgsql2:
    driver: pgsql
    address: localhost:5432
    database: mydatabase
    username: myuser
    password: mypass
    conn:
      maxIdle: 30
      maxOpen: 30
      maxLifeTime: 60000
    slowThreshold: 200 # In millisecond
    skipDefaultTransaction: false
    ignoreRecordNotFoundError: false
  mssql:
    driver: mssql
    address: localhost:5432
    database: mydatabase
    username: myuser
    password: mypass
    conn:
      maxIdle: 30
      maxOpen: 30
      maxLifeTime: 60000

otel:
  trace:
    exporter:
      type: "otlp-grpc"
      otlp:
        grpc:
          address: localhost:4317
          timeout: 5 # Otel collector onnection timeout in seconds
          clientMaxReceiveMessageSize: 4MB # Must have bytesize suffix, such as B, KB, MB, GB
      disabled: false
  metric:
    instrumentationName: "myApp"
    exporter:
      type: "otlp-grpc"
      otlp:
        grpc:
          address: localhost:4317
          timeout: 5 # Otel collector onnection timeout in seconds
          clientMaxReceiveMessageSize: 4MB # Must have bytesize suffix, such as B, KB, MB, GB
      disabled: false
  disabled: false # If true then otel init will be skipped

google:
  credential: XJDKERdlshtkdlkslritjhdsldfjdlerh23943hlk45hlsdkl # Credential JSON key
  cloud:
    pubsub:
      publisher:
        logging:
          payloadLogSizeLimit: 2KB
          stdout: false # Show incoming and outgoing message globally, default is false
          database: pgsql1 # Write log to database with ID pgsql1
      subscriber:
        logging:
          payloadLogSizeLimit: 2KB
          stdout: false # Show incoming and outgoing message globally, default is false
          database: pgsql1 # Write log to database with ID pgsql1
    oauth2:
      verification:
        - aud: 32555940559.apps.googleusercontent.com
          email: playground@gserviceaccount.com
          sub: "114023262557275904507"

zeebe:
  address: 1fcab5e1-1649-4c0b-ab53-cd4375877dc6.syd-1.zeebe.camunda.io:443
  plainTextConn: false
  clientId: j6V~dc0.ypHUKz0SqAu~RpwLf3mKQE7h
  clientSecret: t-MmjpLsq4qWV3lXyPxHTDL6SfHocf3PwVEeh0sTfCS1SK0BuIlwJx4u53NPOyL8
  authorizationServerURL: https://login.cloud.camunda.io/oauth/token

elasticSearch:
  default:
    addresses:
      - localhost
      - 127.0.0.1
    username: myUser
    password: myPass
    logging:
      payloadLogSizeLimit: 2KB
      stdout: false # Show incoming and outgoing message globally, default is false
      db: pgsql1 # Write log to database with ID pgsql1
```

By default it will read config yaml file located in **conf/app.yaml** in **root folder** of the project.

Config file location can be modified with providing its relative path as a **flag argument**
when executing the binary file.

Example:

```shell
my-service-binary --config-path=my-custom-folder/app.yaml
```

