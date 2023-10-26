package pub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/rosaekapratama/go-starter/config"
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
	cfg := config.Instance.GetObject().Google.Cloud.Pubsub.Publisher
	topic := client.Topic(topicId)
	if nil != opts {
		for _, opt := range opts {
			opt.Apply(topic)
		}
	}
	return &PublisherImpl{topic: topic, logging: cfg.Logging}
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

// WithLogging enable or disable outgoing message logging
func (p *PublisherImpl) WithLogging(logging bool) Publisher {
	p.logging = logging
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

	pubsubFields := make(map[string]interface{})
	pubsubFields[constant.LogTypeFieldKey] = constant.LogTypePubsub
	pubsubFields[constant.IsSubscriberFieldKey] = false
	pubsubFields[constant.TopicIdFieldKey] = p.topic.ID()
	pubsubFields[constant.MessageDataFieldKey] = string(message.Data)

	result := p.topic.Publish(ctx, message)
	// Block until the result is returned and
	// a server-generated ID is returned for the published message.
	serverId, err = result.Get(ctx)
	if err != nil {
		pubsubFields[constant.ErrorFieldKey] = err.Error()
		log.WithTraceFields(ctx).WithFields(pubsubFields).GetLogrusLogger().Error()
	} else if p.logging {
		pubsubFields[constant.MessageIdFieldKey] = serverId
		log.WithTraceFields(ctx).WithFields(pubsubFields).GetLogrusLogger().Info()
	}

	return
}

func (p *PublisherImpl) BatchPublish(ctx context.Context, batchData []interface{}, opts ...PublishOption) error {
	log.Tracef(ctx, "Batch publish data, topicId=%s, dataLen=%d", p.topic.ID(), len(batchData))
	errs := cmap.New[bool]()
	wg := &sync.WaitGroup{}
	for i, data := range batchData {
		wg.Add(integer.One)
		traceParent := myContext.TraceParentFromContext(ctx)
		newCtx := myContext.ContextWithTraceParent(context.Background(), traceParent)
		go func(ctx context.Context, i int, wg *sync.WaitGroup, data interface{}, opts ...PublishOption) {
			defer wg.Done()
			_, err := p.Publish(ctx, data, opts...)
			if err != nil {
				errs.Set(strconv.Itoa(i), true)
			}
		}(newCtx, i, wg, data, opts...)
	}
	wg.Wait()

	if errs.Count() > integer.Zero {
		return response.GeneralError
	}
	return nil
}
