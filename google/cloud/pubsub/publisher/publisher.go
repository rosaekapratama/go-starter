package publisher

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/inhies/go-bytesize"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/location"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	myContext "github.com/rosaekapratama/go-starter/context"
	myPubsub "github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
	myOtel "github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/response"
	"strconv"
	"sync"
	"time"
)

const spanPublish = "common.google.cloud.pubsub.Publish %s"

// NewPublisher create new pub
func NewPublisher(client *pubsub.Client, topicId string, opts ...TopicOption) Publisher {
	cfg := config.Instance.GetObject().Google.Cloud.Pubsub.Publisher
	topic := client.Topic(topicId)
	if nil != opts {
		for _, opt := range opts {
			opt.Apply(topic)
		}
	}
	return &publisherImpl{topic: topic, logging: cfg.Logging}
}

// WithJsonEncoder set publisher encoder with JSON encoder
func (p *publisherImpl) WithJsonEncoder() Publisher {
	p.encoder = &jsonEncoder{}
	return p
}

// WithAvroEncoder set publisher encoder with AVRO encoder
func (p *publisherImpl) WithAvroEncoder(schemaName string) Publisher {
	p.encoder = &avroEncoder{schemaName: schemaName}
	return p
}

// WithProtobufEncoder set publisher encoder with protobuf encoder
func (p *publisherImpl) WithProtobufEncoder() Publisher {
	p.encoder = &protobufEncoder{}
	return p
}

// WithStdoutLogging enable or disable outgoing message logging to stdout
func (p *publisherImpl) WithStdoutLogging(logging bool) Publisher {
	p.logging.Stdout = logging
	return p
}

// WithDatabaseLogging enable or disable outgoing message logging to database
func (p *publisherImpl) WithDatabaseLogging(connectionId string) Publisher {
	p.logging.Database = connectionId
	return p
}

// initMessage Init pubsub message
func (p *publisherImpl) initMessage(ctx context.Context, data interface{}, opts ...PublishOption) (*pubsub.Message, error) {
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
	if p.encoder != nil {
		topicConfig, err := p.topic.Config(ctx)
		if err != nil {
			log.Errorf(ctx, err, "[Pub/Sub] Failed to get topic config, topicId=%s", p.topic.ID())
			return nil, err
		}

		bytes, err := p.encoder.Encode(ctx, topicConfig.SchemaSettings, data)
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
func (p *publisherImpl) Publish(ctx context.Context, data interface{}, opts ...PublishOption) (serverId string, err error) {
	ctx, span := myOtel.Trace(ctx, fmt.Sprintf(spanPublish, p.topic.ID()))
	defer span.End()

	message, err := p.initMessage(ctx, data, opts...)
	if err != nil {
		log.Error(ctx, err, "Failed to init pubsub message")
		return str.Empty, err
	}

	result := p.topic.Publish(ctx, message)
	// Block until the result is returned and
	// a server-generated ID is returned for the published message.
	serverId, err = result.Get(ctx)

	if p.logging.Stdout {
		_payloadLogSizeLimit, err := bytesize.Parse(config.Instance.GetObject().Google.Cloud.Pubsub.Publisher.Logging.PayloadLogSizeLimit)
		if err != nil {
			log.Fatal(ctx, err, "Invalid value of GCP pubsub publisher payloadLogSizeLimit config")
		}
		payloadLogSizeLimit := int(_payloadLogSizeLimit)

		pubsubFields := make(map[string]interface{})
		pubsubFields[constant.LogTypeFieldLogKey] = constant.LogTypePubSub
		pubsubFields[constant.IsSubscriberLogKey] = false
		pubsubFields[constant.TopicIdLogKey] = p.topic.ID()
		if len(message.Data) > payloadLogSizeLimit {
			pubsubFields[constant.MessageDataLogKey] = string(message.Data[:payloadLogSizeLimit]) + sym.Ellipsis
		} else if len(message.Data) > 0 {
			pubsubFields[constant.MessageDataLogKey] = string(message.Data)
		} else {
			pubsubFields[constant.MessageDataLogKey] = str.Empty
		}
		if err != nil {
			pubsubFields[constant.ErrorLogKey] = err.Error()
			log.WithTraceFields(ctx).WithFields(pubsubFields).GetLogrusLogger().Error()
		} else if p.logging.Stdout {
			pubsubFields[constant.MessageIdLogKey] = serverId
			log.WithTraceFields(ctx).WithFields(pubsubFields).GetLogrusLogger().Info()
		}
	}

	return
}

func (p *publisherImpl) BatchPublish(ctx context.Context, batchData []interface{}, opts ...PublishOption) error {
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
