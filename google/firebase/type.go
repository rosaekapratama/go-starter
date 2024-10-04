package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type App interface {
	GetMessagingClient() MessagingClient
}

type appImpl struct {
	app             *firebase.App
	messagingClient *messaging.Client
}

type MessagingClient interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
}
