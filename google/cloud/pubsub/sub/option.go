package sub

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"github.com/hamba/avro/v2"
	myAvro "github.com/rosaekapratama/go-starter/avro"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
	"golang.org/x/net/context"
	"reflect"
)

type SubscriptionOption interface {
	Apply(sub *pubsub.Subscription)
}

type Decoder interface {
	Apply(ctx context.Context, message *pubsub.Message) (interface{}, error)
}

type maxOutstandingMessagesOption struct {
	MaxOutstandingMessages int
}

type maxOutstandingBytesOption struct {
	MaxOutstandingBytes int
}

type numGoroutinesOption struct {
	NumGoroutines int
}

type jsonDecoder struct {
	model interface{}
}

type avroDecoder struct {
	schema string
	model  interface{}
}

func (o *maxOutstandingMessagesOption) Apply(sub *pubsub.Subscription) {
	sub.ReceiveSettings.MaxOutstandingMessages = o.MaxOutstandingMessages
}

func (o *maxOutstandingBytesOption) Apply(sub *pubsub.Subscription) {
	sub.ReceiveSettings.MaxOutstandingBytes = o.MaxOutstandingBytes
}

func (o *numGoroutinesOption) Apply(sub *pubsub.Subscription) {
	sub.ReceiveSettings.NumGoroutines = o.NumGoroutines
}

func (o *jsonDecoder) Apply(ctx context.Context, message *pubsub.Message) (interface{}, error) {
	var data interface{}
	if o.model == nil {
		data = make(map[string]interface{})
	} else {
		data = reflect.New(reflect.TypeOf(o.model)).Interface()
	}
	err := json.Unmarshal(message.Data, data)
	if err != nil {
		log.Errorf(ctx, err, "[Pub/Sub] Failed to marshal object, data=%v", message.Data)
		return nil, err
	}
	log.Tracef(ctx, "[Pub/Sub] Prepare JSON-encoded message, data=%#v", data)
	return data, nil
}

func (o *avroDecoder) Apply(ctx context.Context, message *pubsub.Message) (interface{}, error) {
	if o.model == nil {
		log.Errorf(ctx, response.InvalidConfig, "Model must not be null for avro decoding, schemaName=%s", o.schema)
		return nil, response.InvalidConfig
	}

	schema, err := myAvro.SchemaManager.GetSchema(ctx, o.schema)
	if err != nil {
		log.Errorf(ctx, err, "[Pub/Sub] Failed to get avro schema, name=%s", o.schema)
		return nil, err
	}

	data := reflect.New(reflect.TypeOf(o.model)).Interface()
	err = avro.Unmarshal(schema, message.Data, data)
	if err != nil {
		log.Errorf(ctx, err, "[Pub/Sub] codec.TextualFromNative error, data=%v", message.Data)
		return nil, err
	}
	log.Tracef(ctx, "[Pub/Sub] Prepare a JSON-encoded message, data=%#v", data)
	return data, nil
}

func WithMaxOutstandingMessages(maxOutstandingMessages int) SubscriptionOption {
	return &maxOutstandingMessagesOption{MaxOutstandingMessages: maxOutstandingMessages}
}

func WithMaxOutstandingBytes(maxOutstandingBytes int) SubscriptionOption {
	return &maxOutstandingBytesOption{MaxOutstandingBytes: maxOutstandingBytes}
}

func WithNumGoroutines(numGoroutines int) SubscriptionOption {
	return &numGoroutinesOption{NumGoroutines: numGoroutines}
}

func JsonDecoder(model interface{}) Decoder {
	return &jsonDecoder{model: model}
}

func AvroDecoder(schema string, model interface{}) Decoder {
	return &avroDecoder{
		schema: schema,
		model:  model,
	}
}
