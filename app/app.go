package app

import (
	"context"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/rosaekapratama/go-starter/avro"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/database"
	"github.com/rosaekapratama/go-starter/elasticsearch"
	"github.com/rosaekapratama/go-starter/google"
	"github.com/rosaekapratama/go-starter/google/cloud/oauth"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub/subscriber"
	"github.com/rosaekapratama/go-starter/google/cloud/scheduler"
	"github.com/rosaekapratama/go-starter/google/cloud/storage"
	"github.com/rosaekapratama/go-starter/google/drive"
	"github.com/rosaekapratama/go-starter/google/firebase"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/transport/repositories"
	"github.com/rosaekapratama/go-starter/loginit"
	"github.com/rosaekapratama/go-starter/mocks"
	myOtel "github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/redis"
	"github.com/rosaekapratama/go-starter/transport/grpcclient"
	"github.com/rosaekapratama/go-starter/transport/grpcserver"
	"github.com/rosaekapratama/go-starter/transport/restclient"
	"github.com/rosaekapratama/go-starter/transport/restserver"
	"github.com/rosaekapratama/go-starter/transport/soapclient"
	"github.com/rosaekapratama/go-starter/zeebe"
	"github.com/stretchr/testify/mock"
)

var (
	configObject *config.Object
)

func init() {
	ctx := context.Background()

	// Retrieve the main module's information
	info, ok := debug.ReadBuildInfo()
	if ok {
		// Print the version of a specific module (e.g., gin-gonic/gin)
		for _, mod := range info.Deps {
			if mod.Path == "github.com/rosaekapratama/go-starter" {
				log.Infof(ctx, "go-common-gg version: %s", mod.Version)
				break
			}
		}
	}

	args := os.Args
	if strings.HasSuffix(args[0], ".test") {
		initTest(ctx)
	} else {
		initRun(ctx)
	}
}

func initTest(_ context.Context) {
	// To handle init() function which calls config
	mockConfig := mocks.GetMockConfig()
	mockConfig.On("GetString", mock.Anything).Return("string", nil)
	mockConfig.On("GetInt", mock.Anything).Return(integer.Zero, nil)
	mockConfig.On("GetBool", mock.Anything).Return(false, nil)
	mockConfig.On("GetSlice", mock.Anything).Return(make([]interface{}, integer.Zero), nil)
	mockConfig.On("GetStringAndThrowFatalIfEmpty", mock.Anything).Return("string", nil)
	config.Instance = mockConfig
}

func initRun(ctx context.Context) {
	// Init config package
	config.Init()
	configInstance := config.Instance
	configObject = configInstance.GetObject()

	// Init google credential
	projectId := configObject.App.Mode
	credentials, jsonKey := google.CreateCredentials(ctx, configInstance)

	// Extract project ID from credentials
	if credentials != nil && credentials.ProjectID != str.Empty {
		projectId = credentials.ProjectID
	}

	// Set project ID for loginit
	loginit.SetProjectId(projectId)
	log.Infof(ctx, "projectId=%s", projectId)

	// Init google package
	if credentials != nil {
		firebaseApp := firebase.New(ctx, credentials)
		oauthClient := oauth.NewClient(ctx)
		pubsubClient := pubsub.NewClient(ctx, credentials)
		subscriber.Init(pubsubClient)
		scheduler.Init(ctx, credentials)
		schedulerService := scheduler.Service
		storage.Init(ctx, credentials)
		storageClient := storage.Client
		drive.Init(ctx, credentials)
		driveService := drive.Service
		google.Init(
			ctx,
			credentials,
			jsonKey,
			firebaseApp,
			oauthClient,
			pubsubClient,
			schedulerService,
			storageClient,
			driveService,
		)
	}

	// Init log package
	log.Init(ctx, configInstance, projectId)

	// Init otel package
	myOtel.Init(ctx, configInstance)

	// Init avro package
	avro.Init(ctx, configInstance)

	// Init database package
	database.Init(ctx, configInstance)

	// Init redis package
	redis.Init(ctx, configInstance)

	// Init zeebe package
	zeebe.Init(ctx, configInstance)

	// Init server rest logging repository
	var serverRestLogRepository repositories.ITransportLogRepository
	dbLogRestServer := configInstance.GetObject().Transport.Server.Rest.Logging.Database
	if dbLogRestServer != str.Empty {
		DB, _, err := database.Manager.DB(ctx, dbLogRestServer)
		if err != nil {
			log.Fatalf(ctx, err, "Failed to find database ID '%s'", dbLogRestServer)
			return
		}
		serverRestLogRepository = repositories.NewTransportLogRepository(DB)
	}

	// Init client rest logging repository
	var clientRestLogRepository repositories.ITransportLogRepository
	dbLogRestClient := configInstance.GetObject().Transport.Client.Rest.Logging.Database
	if dbLogRestClient != str.Empty {
		DB, _, err := database.Manager.DB(ctx, dbLogRestClient)
		if err != nil {
			log.Fatalf(ctx, err, "Failed to find database ID '%s'", dbLogRestClient)
			return
		}
		clientRestLogRepository = repositories.NewTransportLogRepository(DB)
	}

	// Init client soap logging repository
	var clientSoapLogRepository repositories.ITransportLogRepository
	dbLogSoapClient := configInstance.GetObject().Transport.Client.Soap.Logging.Database
	if dbLogSoapClient != str.Empty {
		DB, _, err := database.Manager.DB(ctx, dbLogSoapClient)
		if err != nil {
			log.Fatalf(ctx, err, "Failed to find database ID '%s'", dbLogSoapClient)
			return
		}
		clientSoapLogRepository = repositories.NewTransportLogRepository(DB)
	}

	// Init SOAP client
	soapclient.Init(ctx, configInstance, clientSoapLogRepository)

	// Init REST client
	restclient.Init(ctx, configInstance, clientRestLogRepository)

	// Init REST server
	restserver.Init(ctx, configInstance, serverRestLogRepository)

	// Init GRPC server
	grpcserver.Init(ctx, configInstance)

	// Init GRPC client
	grpcclient.Init(ctx, configInstance)

	// Init elasticsearch package
	elasticsearch.Init(ctx, configInstance)
}

func Run() {
	// Run REST server
	go restserver.Run()

	// Run GRPC server
	go grpcserver.Run()

	group := sync.WaitGroup{}
	group.Add(1)
	group.Wait()
}
