package app

import (
	"context"
	"github.com/rosaekapratama/go-starter/avro"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/database"
	"github.com/rosaekapratama/go-starter/elasticsearch"
	"github.com/rosaekapratama/go-starter/google"
	"github.com/rosaekapratama/go-starter/google/cloud/oauth"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub/sub"
	"github.com/rosaekapratama/go-starter/google/cloud/scheduler"
	"github.com/rosaekapratama/go-starter/google/cloud/storage"
	"github.com/rosaekapratama/go-starter/google/firebase"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/loginit"
	"github.com/rosaekapratama/go-starter/mocks"
	myOtel "github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/redis"
	"github.com/rosaekapratama/go-starter/transport/grpcserver"
	"github.com/rosaekapratama/go-starter/transport/logging/repositories"
	"github.com/rosaekapratama/go-starter/transport/restclient"
	"github.com/rosaekapratama/go-starter/transport/restserver"
	"github.com/rosaekapratama/go-starter/zeebe"
	"github.com/stretchr/testify/mock"
	"os"
	"strings"
	"sync"
)

var (
	configObject *config.Object
)

func init() {
	ctx := context.Background()

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
		oauthClient := oauth.New(ctx)
		pubsubClient := pubsub.New(ctx, credentials)
		sub.Init(pubsubClient)
		scheduler.Init(ctx, credentials)
		schedulerService := scheduler.Service
		storage.Init(ctx, credentials)
		storageClient := storage.Client
		google.Init(
			ctx,
			credentials,
			jsonKey,
			firebaseApp,
			oauthClient,
			pubsubClient,
			schedulerService,
			storageClient)
	}

	// Init log package
	log.Init(ctx, configInstance, projectId)

	// Init otel package
	myOtel.Init(ctx, configInstance)

	// Init avro package
	avro.Init(ctx, configInstance)

	// Init database package
	database.Init(ctx, configInstance)

	// Init elasticsearch package
	elasticsearch.Init(ctx, configInstance)

	// Init redis package
	redis.Init(ctx, configInstance)

	// Init zeebe package
	zeebe.Init(ctx, configInstance)

	// Init rest logging repository
	var restLogRepository repositories.IRestLogRepository
	dbLog := configInstance.GetObject().Transport.Server.Rest.Logging.Database
	if dbLog != str.Empty {
		DB, _, err := database.Manager.DB(ctx, dbLog)
		if err != nil {
			log.Fatalf(ctx, err, "Failed to find database ID '%s'", dbLog)
			return
		}
		restLogRepository = repositories.NewRestLogRepository(DB)
	}

	// Init http client package
	restclient.Init(ctx, configInstance, restLogRepository)

	// Init REST server
	restserver.Init(ctx, configInstance, restLogRepository)

	// Init GRPC server
	grpcserver.Init(ctx, configInstance)
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
