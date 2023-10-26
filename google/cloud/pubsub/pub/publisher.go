package pub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/location"
	"github.com/rosaekapratama/go-starter/constant/str"
	myContext "github.com/rosaekapratama/go-starter/context"
	myPubsub "github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub/constant"
	"github.com/rosaekapratama/go-starter/log"
	myOtel "github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/response"
	"strconv"
	"sync"
	"time"
)

const spanPublish = "pubsub publish %s"

// NewPublisher create new pub
func NewPublisher(client *pubsub.Client, topicId string, opts ...TopicOption) Publisher {
	topic := client.Topic(topicId)
	if nil != opts {
		for _, opt := range opts {
			opt.Apply(topic)
		}
	}
	return &PublisherImpl{topic: topic}
}

// WithJsonEncoder set publisher encoder with JSON encoder
func (p *PublisherImpl) WithJsonEncoder() Publisher {
	p.encoder = &jsonEncoder{}
	return p
}

// WithAvroEncoder set publisher encoder with AVRO encoder
func (p *PublisherImpl) WithAvroEncoder(schemaName string) Publisher {
	p.encoder = &avroEncoder{schemaName: schemaName}
	return p
}

// initMessage Init pubsub message
func initMessage(ctx context.Context, data interface{}, encoder Encoder, opts ...PublishOption) (*pubsub.Message, error) {
	message := &pubsub.Message{
		Attributes: map[string]string{
			myPubsub.TraceparentAttrKey:       myContext.TraceParentFromContext(ctx),
			myPubsub.OriginPublishTimeAttrKey: time.Now().In(location.AsiaJakarta).Format(time.RFC3339),
		},
	}

	// Apply options
	if nil != opts {
		for _, opt := range opts {
			opt.Apply(message)
		}
	}

	// Encode message
	if encoder != nil {
		bytes, err := encoder.Encode(ctx, data)
		if err != nil {
			log.Error(ctx, err, "[Pub/Sub] Failed to encode data")
			return nil, err
		}
		message.Data = bytes
	} else {
		switch data := data.(type) {
		case string:
			message.Data = []byte(data)
		case int:
			message.Data = []byte(strconv.Itoa(data))
		case bool:
			message.Data = []byte(strconv.FormatBool(data))
		default:
			log.Errorf(ctx, response.UnsupportedType, "[Pub/Sub] Unsupported type on publish, type=%T", data)
			return nil, response.UnsupportedType
		}
	}
	return message, nil
}

// Publish publishing data to pubsub topic
func (p *PublisherImpl) Publish(ctx context.Context, data interface{}, opts ...PublishOption) (serverId string, err error) {
	ctx, span := myOtel.Trace(ctx, fmt.Sprintf(spanPublish, p.topic.ID()))
	defer span.End()

	message, err := initMessage(ctx, data, p.encoder, opts...)
	if err != nil {
		log.Error(ctx, err, "Failed to init pubsub message")
		return str.Empty, err
	}

	result := p.topic.Publish(ctx, message)
	// Block until the result is returned and
	// a server-generated ID is returned for the published message.
	serverId, err = result.Get(ctx)
	if err != nil {
		log.Error(ctx, err)
	}

	pubsubFields := make(map[string]interface{})
	pubsubFields[constant.LogTypeFieldKey] = constant.LogTypePubsub
	pubsubFields[constant.IsSubscriberFieldKey] = false
	pubsubFields[constant.TopicIdFieldKey] = p.topic.ID()
	pubsubFields[constant.MessageIdFieldKey] = serverId
	pubsubFields[constant.MessageDataFieldKey] = message
	log.WithTraceFields(ctx).WithFields(pubsubFields).GetLogrusLogger().Info()

	return
}

func batchPublish(ctx context.Context, wg *sync.WaitGroup, topic *pubsub.Topic, data interface{}, encoder Encoder, errs []error, opts ...PublishOption) []error {
	ctx, span := myOtel.Trace(ctx, fmt.Sprintf(spanPublish, topic.ID()))
	defer span.End()

	wg.Add(integer.One)
	message, err := initMessage(ctx, data, encoder, opts...)
	if err != nil {
		log.Error(ctx, err, "Failed to init pubsub message, topicId=%s", topic.ID())
		return []error{err}
	}
	result := topic.Publish(ctx, message)
	log.Tracef(ctx, "[Pub/Sub] Publish message, topicId=%s", topic.ID())
	go func(res *pubsub.PublishResult) {
		defer wg.Done()
		// The Get method blocks until a server-generated ID or
		// an error is returned for the published message.
		_, err := res.Get(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}(result)

	return errs
}

func (p *PublisherImpl) BatchPublish(ctx context.Context, batchData []interface{}, opts ...PublishOption) error {
	log.Tracef(ctx, "Batch publish data, topicId=%s, dataLen=%d", p.topic.ID(), len(batchData))
	errs := make([]error, integer.Zero)
	var wg *sync.WaitGroup
	for _, data := range batchData {
		errs = batchPublish(ctx, wg, p.topic, data, p.encoder, errs, opts...)
	}
	wg.Wait()

	if len(errs) > integer.Zero {
		for _, err := range errs {
			log.Error(ctx, err)
		}
		return response.GeneralError
	}
	return nil
}
