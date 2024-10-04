package google

import (
	"cloud.google.com/go/pubsub"
	"github.com/rosaekapratama/go-starter/google/cloud/oauth"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub/publisher"
	"github.com/rosaekapratama/go-starter/google/cloud/scheduler"
	"github.com/rosaekapratama/go-starter/google/cloud/storage"
	"github.com/rosaekapratama/go-starter/google/drive"
	"github.com/rosaekapratama/go-starter/google/firebase"
	"golang.org/x/oauth2/google"
)

type IManager interface {
	GetJsonKey() *JsonKey
	GetCredentials() *google.Credentials
	GetFirebaseApp() firebase.App
	GetOAuthClient() oauth.Client
	GetPubSubClient() *pubsub.Client
	NewPubSubPublisher(topicId string, opts ...publisher.TopicOption) publisher.Publisher
	GetSchedulerService() scheduler.IService
	GetStorageClient() storage.IClient
	GetDriveService() drive.IService
}

type managerImpl struct {
	credentials *google.Credentials
	jsonKey     *JsonKey

	firebaseApp      firebase.App
	oauthClient      oauth.Client
	pubsubClient     *pubsub.Client
	schedulerService scheduler.IService
	storageClient    storage.IClient
	driveService     drive.IService
}

type JsonKey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}
