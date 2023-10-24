package pub

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"github.com/hamba/avro/v2"
	myAvro "github.com/rosaekapratama/go-starter/avro"
	pubsub2 "github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/log"
	"golang.org/x/net/context"
	"time"
)

var (
	schemaManager = myAvro.SchemaManager
)

type Encoder interface {
	Encode(ctx context.Context, data interface{}) ([]byte, error)
}

type jsonEncoder struct {
}

func (e *jsonEncoder) Encode(_ context.Context, data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

type avroEncoder struct {
	schemaName string
}

func (e *avroEncoder) Encode(ctx context.Context, data interface{}) ([]byte, error) {
	schema, err := schemaManager.GetSchema(ctx, e.schemaName)
	if err != nil {
		log.Errorf(ctx, err, "[Pub/Sub] Failed to get avro schema, schemaName=%s", e.schemaName)
		return nil, err
	}
	bytes, err := avro.Marshal(schema, data)
	if err != nil {
		log.Error(ctx, err, "[Pub/Sub] avro.Marshal error")
		return nil, err
	}
	log.Trace(ctx, "[Pub/Sub] Prepare avro-encoded message")
	return bytes, nil
}

type PublishOption interface {
	Apply(message *pubsub.Message)
}

type attrOption struct {
	key   string
	value string
}

type stateOption struct {
	state pubsub2.State
}

func (o *attrOption) Apply(message *pubsub.Message) {
	message.Attributes[o.key] = o.value
}

func (o *stateOption) Apply(message *pubsub.Message) {
	message.Attributes[pubsub2.StateAttrKey] = string(o.state)
}

func WithAttribute(key string, value string) PublishOption {
	return &attrOption{
		key:   key,
		value: value,
	}
}

func WithState(state pubsub2.State) PublishOption {
	return &stateOption{
		state: state,
	}
}

type TopicOption interface {
	Apply(*pubsub.Topic)
}

type byteThresholdOption struct {
	byteThreshold int
}

type countThresholdOption struct {
	countThreshold int
}

type delayThresholdOption struct {
	delayThreshold time.Duration
}

type timeoutOption struct {
	timeout time.Duration
}

func (o *byteThresholdOption) Apply(topic *pubsub.Topic) {
	topic.PublishSettings.ByteThreshold = o.byteThreshold
}

func (o *countThresholdOption) Apply(topic *pubsub.Topic) {
	topic.PublishSettings.CountThreshold = o.countThreshold
}

func (o *delayThresholdOption) Apply(topic *pubsub.Topic) {
	topic.PublishSettings.DelayThreshold = o.delayThreshold
}

func (o *timeoutOption) Apply(topic *pubsub.Topic) {
	topic.PublishSettings.Timeout = o.timeout
}

func WithByteThreshold(byteThreshold int) TopicOption {
	return &byteThresholdOption{byteThreshold: byteThreshold}
}

func WithCountThreshold(countThreshold int) TopicOption {
	return &countThresholdOption{countThreshold: countThreshold}
}

func WithDelayThreshold(delayThreshold time.Duration) TopicOption {
	return &delayThresholdOption{delayThreshold: delayThreshold}
}

func WithTimeout(timeout time.Duration) TopicOption {
	return &timeoutOption{timeout: timeout}
}
