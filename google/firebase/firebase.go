package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/rosaekapratama/go-starter/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func New(ctx context.Context, credentials *google.Credentials) App {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentials(credentials))
	if err != nil {
		log.Fatal(ctx, err, "Failed to init firebase app")
		return nil
	}
	log.Info(ctx, "Google firebase app service is initiated")

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatal(ctx, err, "Failed to init firebase cloud messaging client")
		return nil
	}
	log.Info(ctx, "Google firebase cloud messaging client service is initiated")

	return &appImpl{
		app:             app,
		messagingClient: messagingClient,
	}
}

func (app *appImpl) GetMessagingClient() MessagingClient {
	return app.messagingClient
}
