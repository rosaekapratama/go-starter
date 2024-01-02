package subscriber

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"fmt"
	"github.com/hamba/avro/v2"
	myAvro "github.com/rosaekapratama/go-starter/avro"
	"github.com/rosaekapratama/go-starter/log"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type SubscriptionOption interface {
	// ApplyMessageOption is for message option
	ApplyMessageOption(sub *pubsub.Subscription)

	// ApplyDecoderOption is for decoder option
	ApplyDecoderOption(ctx context.Context, message *pubsub.Message) (interface{}, error)
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

type jsonDecoderOption[T any] struct {
}

type avroDecoderOption[T any] struct {
	schema string
}

type protobufDecoderOption[T any] struct {
}

func (o *maxOutstandingMessagesOption) ApplyMessageOption(sub *pubsub.Subscription) {
	sub.ReceiveSettings.MaxOutstandingMessages = o.MaxOutstandingMessages
}

func (o *maxOutstandingMessagesOption) ApplyDecoderOption(_ context.Context, _ *pubsub.Message) (v interface{}, err error) {
	return
}

func (o *maxOutstandingBytesOption) ApplyMessageOption(sub *pubsub.Subscription) {
	sub.ReceiveSettings.MaxOutstandingBytes = o.MaxOutstandingBytes
}

func (o *maxOutstandingBytesOption) ApplyDecoderOption(_ context.Context, _ *pubsub.Message) (v interface{}, err error) {
	return
}

func (o *numGoroutinesOption) ApplyMessageOption(sub *pubsub.Subscription) {
	sub.ReceiveSettings.NumGoroutines = o.NumGoroutines
}

func (o *numGoroutinesOption) ApplyDecoderOption(_ context.Context, _ *pubsub.Message) (v interface{}, err error) {
	return
}

func (o *jsonDecoderOption[T]) ApplyMessageOption(_ *pubsub.Subscription) {
}

func (o *jsonDecoderOption[T]) ApplyDecoderOption(ctx context.Context, message *pubsub.Message) (interface{}, error) {
	obj := new(T)
	err := json.Unmarshal(message.Data, obj)
	if err != nil {
		log.Errorf(ctx, err, "[Pub/Sub] Failed to marshal object, data=%v", message.Data)
		return nil, err
	}
	return obj, nil
}

func (o *avroDecoderOption[T]) ApplyMessageOption(_ *pubsub.Subscription) {
}

func (o *avroDecoderOption[T]) ApplyDecoderOption(ctx context.Context, message *pubsub.Message) (interface{}, error) {
	schema, err := myAvro.SchemaManager.GetSchema(ctx, o.schema)
	if err != nil {
		log.Errorf(ctx, err, "[Pub/Sub] Failed to get avro schema, name=%s", o.schema)
		return nil, err
	}

	obj := new(T)
	err = avro.Unmarshal(schema, message.Data, obj)
	if err != nil {
		log.Errorf(ctx, err, "[Pub/Sub] codec.TextualFromNative error, data=%v", message.Data)
		return nil, err
	}
	return obj, nil
}

func (o *protobufDecoderOption[T]) ApplyMessageOption(_ *pubsub.Subscription) {
}

func (o *protobufDecoderOption[T]) ApplyDecoderOption(ctx context.Context, message *pubsub.Message) (obj interface{}, err error) {
	obj = new(T)
	if _, ok := obj.(proto.Message); !ok {
		return nil, fmt.Errorf("pubsub message object must implement proto.Message, type=%T", obj)
	}

	encoding := message.Attributes["googclient_schemaencoding"]
	switch encoding {
	case "BINARY":
		if err = proto.Unmarshal(message.Data, obj.(proto.Message)); err != nil {
			log.Error(ctx, err, "[Pub/Sub] proto.Unmarshal error")
			return
		}
	case "JSON":
		if err = protojson.Unmarshal(message.Data, obj.(proto.Message)); err != nil {
			log.Error(ctx, err, "[Pub/Sub] protojson.Unmarshal error")
			return
		}
	default:
		err = fmt.Errorf("invalid subscriber protobuf encoding type, type=%v", encoding)
		log.Errorf(ctx, err, "[Pub/Sub] Invalid subscriber encoding type, type=%v", encoding)
	}
	return
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

func WithJsonDecoder[T any]() SubscriptionOption {
	return &jsonDecoderOption[T]{}
}

func WithAvroDecoder[T any](schema string) SubscriptionOption {
	return &avroDecoderOption[T]{
		schema: schema,
	}
}

func WithProtobufDecoder[T any]() SubscriptionOption {
	return &protobufDecoderOption[T]{}
}
