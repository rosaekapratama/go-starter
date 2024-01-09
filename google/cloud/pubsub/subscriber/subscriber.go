package subscriber

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	myContext "github.com/rosaekapratama/go-starter/context"
	myPubsub "github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
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
func Receive(subId string, f func(ctx context.Context, plainMessage *pubsub.Message, decodedMessage interface{}), opts ...SubscriptionOption) {
	wg.Add(integer.One)
	go receive(subId, myPubsub.StateAny, f, opts...)
}

func ReceiveWithState(subId string, state myPubsub.State, f func(ctx context.Context, plainMessage *pubsub.Message, decodedMessage interface{}), opts ...SubscriptionOption) {
	wg.Add(integer.One)
	go receive(subId, state, f, opts...)
}

func receive(subId string, state myPubsub.State, f func(ctx context.Context, plainMessage *pubsub.Message, decodedMessage interface{}), opts ...SubscriptionOption) {
	ctx := context.Background()
	cfg := config.Instance.GetObject().Google.Cloud.Pubsub.Subscriber

	// Init subscription and apply subscription options
	sub := client.Subscription(subId)
	for _, opt := range opts {
		if opt != nil {
			opt.ApplyMessageOption(sub)
		}
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
	err = sub.Receive(ctx, func(ctx context.Context, plainMessage *pubsub.Message) {
		// Add traceparent to pubsub message attributes
		ctx = myContext.ContextWithTraceParent(ctx, plainMessage.Attributes[myPubsub.TraceparentAttrKey])
		ctx, span := otel.Trace(ctx, fmt.Sprintf(spanSubscriber, subId))
		defer span.End()

		// If state matches, then continue, or break if not matches
		var messageState string
		if v, ok := plainMessage.Attributes[myPubsub.StateAttrKey]; ok {
			messageState = v
		}

		if cfg.Logging.Stdout {
			// Logging incoming message
			pubsubFields := make(map[string]interface{})
			pubsubFields[constant.LogTypeFieldKey] = constant.LogTypePubSub
			pubsubFields[constant.IsSubscriberFieldKey] = true
			pubsubFields[constant.SubscriberIdFieldKey] = subId
			pubsubFields[constant.MessageIdFieldKey] = plainMessage.ID
			pubsubFields[constant.MessageStateFieldKey] = messageState
			if len(plainMessage.Data) > integer.Zero && len(plainMessage.Data) <= (64*1000) {
				pubsubFields[constant.MessageDataFieldKey] = string(plainMessage.Data)
			}
			log.WithTraceFields(ctx).WithFields(pubsubFields).GetLogrusLogger().Info()
		}

		if state != myPubsub.StateAny && string(state) != strings.ToLower(messageState) {
			// Break cause not match
			log.Tracef(ctx, "State doesn't match, subId=%s, subState=%s, msgState=%s", subId, state, messageState)
			plainMessage.Ack()
			return
		}

		// Apply first found decoder
		var decodedMessage interface{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			decodedMessage, err = opt.ApplyDecoderOption(ctx, plainMessage)
			if err != nil {
				log.Errorf(ctx, err, "Failed to apply receiver decoder option, subId=%s, option=%T", subId, opt)
				plainMessage.Nack()
				return
			}
			if decodedMessage != nil {
				break
			}
		}

		if decodedMessage != nil {
			f(ctx, plainMessage, decodedMessage)
		} else {
			f(ctx, plainMessage, nil)
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
