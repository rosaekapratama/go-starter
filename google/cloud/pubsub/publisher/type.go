package publisher

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rosaekapratama/go-starter/config"
)

type Publisher interface {
	Publish(ctx context.Context, data interface{}, opts ...PublishOption) (serverId string, err error)
	BatchPublish(ctx context.Context, batchData []interface{}, opts ...PublishOption) error
	WithJsonEncoder() Publisher
	WithAvroEncoder(schemaName string) Publisher
	WithProtobufEncoder() Publisher
	WithStdoutLogging(logging bool) Publisher
	WithDatabaseLogging(connectionId string) Publisher
}

type publisherImpl struct {
	topic   *pubsub.Topic
	encoder Encoder
	logging *config.GoogleCloudPubsubPublisherLoggingConfig
}
