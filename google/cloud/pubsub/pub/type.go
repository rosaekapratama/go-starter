package pub

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, data interface{}, opts ...PublishOption) (serverId string, err error)
	BatchPublish(ctx context.Context, batchData []interface{}, opts ...PublishOption) error
	WithJsonEncoder() Publisher
	WithAvroEncoder(schemaName string) Publisher
	WithLogging(logging bool) Publisher
}

type PublisherImpl struct {
	topic   *pubsub.Topic
	encoder Encoder
	logging bool
}
