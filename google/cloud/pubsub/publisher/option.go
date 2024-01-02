package publisher

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"fmt"
	"github.com/hamba/avro/v2"
	myAvro "github.com/rosaekapratama/go-starter/avro"
	myPubsub "github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/log"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"time"
)

var (
	schemaManager = myAvro.SchemaManager
)

type Encoder interface {
	Encode(ctx context.Context, schemaSettings *pubsub.SchemaSettings, data interface{}) ([]byte, error)
}

type jsonEncoder struct {
}

func (e *jsonEncoder) Encode(_ context.Context, _ *pubsub.SchemaSettings, data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

type avroEncoder struct {
	schemaName string
}

type protobufEncoder struct {
}

func (e *avroEncoder) Encode(ctx context.Context, _ *pubsub.SchemaSettings, data interface{}) ([]byte, error) {
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
	return bytes, nil
}

func (e *protobufEncoder) Encode(ctx context.Context, schemaSettings *pubsub.SchemaSettings, data interface{}) (bytes []byte, err error) {
	var protoMessage proto.Message
	if v, ok := data.(proto.Message); ok {
		protoMessage = v
	} else {
		return nil, fmt.Errorf("pubsub message object must implement proto.Message, type=%T", data)
	}

	encoding := schemaSettings.Encoding
	switch encoding {
	case pubsub.EncodingBinary:
		bytes, err = proto.Marshal(protoMessage)
		if err != nil {
			log.Error(ctx, err, "[Pub/Sub] proto.Marshal error")
		}
	case pubsub.EncodingJSON:
		bytes, err = protojson.Marshal(protoMessage)
		if err != nil {
			log.Error(ctx, err, "[Pub/Sub] protojson.Marshal error")
		}
	default:
		err = fmt.Errorf("invalid publisher protobuf encoding type, type=%v", encoding)
		log.Errorf(ctx, err, "[Pub/Sub] Invalid publisher encoding type, type=%v", encoding)
	}

	return
}

type PublishOption interface {
	Apply(message *pubsub.Message)
}

type attrOption struct {
	key   string
	value string
}

type stateOption struct {
	state myPubsub.State
}

func (o *attrOption) Apply(message *pubsub.Message) {
	message.Attributes[o.key] = o.value
}

func (o *stateOption) Apply(message *pubsub.Message) {
	message.Attributes[myPubsub.StateAttrKey] = string(o.state)
}

func WithAttribute(key string, value string) PublishOption {
	return &attrOption{
		key:   key,
		value: value,
	}
}

func WithState(state myPubsub.State) PublishOption {
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
