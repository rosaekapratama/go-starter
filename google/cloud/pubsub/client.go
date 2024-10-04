package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rosaekapratama/go-starter/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, credentials *google.Credentials) (client *pubsub.Client) {
	var err error
	client, err = pubsub.NewClient(ctx, credentials.ProjectID, option.WithCredentials(credentials))
	if err != nil {
		log.Fatal(ctx, err, "Failed to create google pubsub client")
	}
	log.Info(ctx, "Google pub/sub client service is initiated")
	return
}
