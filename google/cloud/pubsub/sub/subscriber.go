package sub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	myContext "github.com/rosaekapratama/go-starter/context"
	myPubsub "github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub/constant"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/response"
	"strings"
	"sync"
	"time"
)

const spanSubscriber = "pubsub receive %s"

var (
	client *pubsub.Client
	wg     = sync.WaitGroup{}
)

func Init(newClient *pubsub.Client) {
	client = newClient
}

type Subscriber struct {
}

// Receive
// Example to use
//
//	pubsub.Receive("policy-member-endorse-sub", Test)
//
// where Test is a function with appropriate parameters
//
//	func Test(ctx context.Context, msg *pubsub.Message, data interface{}) {
//		ctx, span := otel.Trace(ctx, "Test")
//		defer span.End()
//		defer msg.Ack()
//		log.Tracef(ctx, msg.ID, string(data))
//	}
func Receive(subId string, f func(ctx context.Context, message *pubsub.Message, data interface{}), decoder Decoder, opts ...SubscriptionOption) {
	wg.Add(integer.One)
	go receive(subId, myPubsub.StateAny, f, decoder, opts...)
}

func ReceiveWithState(subId string, state myPubsub.State, f func(ctx context.Context, message *pubsub.Message, data interface{}), decoder Decoder, opts ...SubscriptionOption) {
	wg.Add(integer.One)
	go receive(subId, state, f, decoder, opts...)
}

func receive(subId string, state myPubsub.State, f func(ctx context.Context, message *pubsub.Message, data interface{}), decoder Decoder, opts ...SubscriptionOption) {
	ctx := context.Background()
	cfg := config.Instance.GetObject().Google.Cloud.Pubsub.Subscriber

	// Init subscription and apply subscription options
	sub := client.Subscription(subId)
	for _, opt := range opts {
		opt.Apply(sub)
	}

	// Check if exists or not
	exists, err := sub.Exists(ctx)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to check subscriber existence, subId=%s", subId)
		return
	} else if !exists {
		log.Fatalf(ctx, response.InitFailed, "Subscriber doesn't exists, subId=%s", subId)
		return
	} else {
		wg.Done()
	}

	wg.Wait()
	log.Infof(ctx, "Start google pubsub sub, subId=%s", subId)
	// Run subscription.Receive function to receive data from pubsub
	err = sub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		// Add traceparent to pubsub message attributes
		ctx = myContext.ContextWithTraceParent(ctx, message.Attributes[myPubsub.TraceparentAttrKey])
		ctx, span := otel.Trace(ctx, fmt.Sprintf(spanSubscriber, subId))
		defer span.End()

		// If state matches, then continue, or break if not matches
		var messageState string
		if v, ok := message.Attributes[myPubsub.StateAttrKey]; ok {
			messageState = v
		}

		if cfg.Logging {
			// Logging incoming message
			pubsubFields := make(map[string]interface{})
			pubsubFields[constant.LogTypeFieldKey] = constant.LogTypePubsub
			pubsubFields[constant.IsSubscriberFieldKey] = true
			pubsubFields[constant.SubscriberIdFieldKey] = subId
			pubsubFields[constant.MessageIdFieldKey] = message.ID
			pubsubFields[constant.MessageStateFieldKey] = messageState
			pubsubFields[constant.MessageDataFieldKey] = string(message.Data)
			log.WithTraceFields(ctx).WithFields(pubsubFields).GetLogrusLogger().Info()
		}

		if state != myPubsub.StateAny && string(state) != strings.ToLower(messageState) {
			// Break cause not match
			log.Tracef(ctx, "State doesn't match, subId=%s, subState=%s, msgState=%s", subId, state, messageState)
			message.Ack()
			return
		}

		// Decode message if decoder not nil
		if decoder != nil {
			data, err := decoder.Apply(ctx, message)
			if err != nil {
				log.Errorf(ctx, err, "Failed to apply receive decoder option, subId=%s, option=%T", subId, decoder)
				return
			}
			f(ctx, message, data)
		} else {
			f(ctx, message, nil)
		}
	})
	if err != nil {
		log.Fatalf(ctx, err, "Failed to init pubsub receiver, subId=%s", subId)
		return
	}
}

// GetOriginPublishTimeFromMessage return nil if publish time attribute not found
func GetOriginPublishTimeFromMessage(ctx context.Context, msg *pubsub.Message) (*time.Time, error) {
	if v, ok := msg.Attributes[myPubsub.OriginPublishTimeAttrKey]; ok {
		publishTime, err := time.Parse(time.RFC3339, v)
		if err != nil {
			log.Errorf(ctx, err, "Failed to parse origin publish time string, publishTime=%s", v)
			return nil, err
		}
		return &publishTime, nil
	}
	return nil, nil
}
