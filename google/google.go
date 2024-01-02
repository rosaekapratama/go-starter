package google

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/secretmanager/apiv1"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/google/cloud/oauth"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub/publisher"
	"github.com/rosaekapratama/go-starter/google/cloud/scheduler"
	"github.com/rosaekapratama/go-starter/google/cloud/storage"
	"github.com/rosaekapratama/go-starter/google/firebase"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
	"golang.org/x/oauth2/google"
	"os"
	"strings"
)

const (
	googleAppCredEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"
)

var (
	errFailedToDecodeGoogleJsonKey            = "failed to decode google json key"
	errFailedToUnmarshalJsonKey               = "failed to unmarshal json key to its struct"
	errFailedToGetGoogleCredentialFromJsonKey = "failed to get google credential from json key"
	errMissingGoogleConfiguration             = "missing google configuration"

	Manager IManager
)

func CreateCredentials(ctx context.Context, config config.Config) (credentials *google.Credentials, jsonKey *JsonKey) {
	cfg := config.GetObject().Google

	if cfg == nil || cfg.Disabled {
		log.Warn(ctx, "Google manager is disabled")
		return
	}

	validateCfg(ctx, cfg)

	var b []byte
	var err error
	if _, exists := os.LookupEnv(googleAppCredEnvVar); !exists {
		// Check if key is json in base64 or not
		if strings.HasPrefix(cfg.Credential, sym.OpenCurlyBracket) &&
			strings.HasSuffix(cfg.Credential, sym.CloseCurlyBracket) {
			b = []byte(cfg.Credential)
		} else {
			b, err = base64.StdEncoding.DecodeString(cfg.Credential)
		}
		if err != nil {
			log.Fatal(ctx, err, errFailedToDecodeGoogleJsonKey)
			return
		}

		// unmarshall to struct for future usage
		jsonKey = &JsonKey{}
		err = json.Unmarshal(b, jsonKey)
		if err != nil {
			log.Fatalf(ctx, err, errFailedToUnmarshalJsonKey)
			return
		}
	}

	// Get google credentials
	if b != nil && len(b) > integer.Zero {
		credentials, err = google.CredentialsFromJSON(ctx, b, secretmanager.DefaultAuthScopes()...)
	} else {
		credentials, err = google.FindDefaultCredentials(ctx, secretmanager.DefaultAuthScopes()...)
	}
	if err != nil {
		log.Fatal(ctx, err, errFailedToGetGoogleCredentialFromJsonKey)
		return
	}

	log.Info(ctx, "Google credentials is initiated")
	return
}

func validateCfg(ctx context.Context, cfg *config.GoogleConfig) {

	// If credential not found return error
	if cfg.Credential == str.Empty {
		log.Fatal(ctx, response.InitFailed, errMissingGoogleConfiguration)
		return
	}
}

func Init(
	ctx context.Context,
	credentials *google.Credentials,
	jsonKey *JsonKey,
	firebaseApp firebase.App,
	oauthClient oauth.Client,
	pubsubClient *pubsub.Client,
	schedulerService scheduler.IService,
	storageClient storage.IClient) {
	Manager = &ManagerImpl{
		credentials:      credentials,
		jsonKey:          jsonKey,
		firebaseApp:      firebaseApp,
		oauthClient:      oauthClient,
		pubsubClient:     pubsubClient,
		schedulerService: schedulerService,
		storageClient:    storageClient,
	}
	log.Info(ctx, "Google manager is initiated")
}

func (m *ManagerImpl) GetJsonKey() *JsonKey {
	return m.jsonKey
}

func (m *ManagerImpl) GetCredentials() *google.Credentials {
	return m.credentials
}

func (m *ManagerImpl) GetFirebaseApp() firebase.App {
	return m.firebaseApp
}

func (m *ManagerImpl) GetOAuthClient() oauth.Client {
	return m.oauthClient
}

func (m *ManagerImpl) GetPubSubClient() *pubsub.Client {
	return m.pubsubClient
}

func (m *ManagerImpl) NewPubSubPublisher(topicId string, opts ...publisher.TopicOption) publisher.Publisher {
	return publisher.NewPublisher(m.pubsubClient, topicId, opts...)
}

func (m *ManagerImpl) GetSchedulerService() scheduler.IService {
	return m.schedulerService
}

func (m *ManagerImpl) GetStorageClient() storage.IClient {
	return m.storageClient
}
