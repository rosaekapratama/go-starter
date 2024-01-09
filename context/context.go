package context

import (
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/log"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	traceparentFormat = "00-%s-%s-01"
	tokenKey          = "token"
	userIdKey         = "userId"
	usernameKey       = "username"
	fullnameKey       = "fullname"
	realmKey          = "realm"
	emailKey          = "email"
)

var (
	keys = []string{tokenKey, userIdKey, usernameKey, fullnameKey, realmKey, emailKey}
)

func GetAllManagedKey() []string {
	return keys
}

func AddManagedKey(key string) {
	keys = append(keys, key)
}

func NewContextFromTraceParent(ctx context.Context) context.Context {
	return ContextWithTraceParent(context.Background(), TraceParentFromContext(ctx))
}

func TraceParentFromContext(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	return fmt.Sprintf(traceparentFormat, sc.TraceID().String(), sc.SpanID().String())
}

func ContextWithTraceParent(parent context.Context, traceparent string) context.Context {
	traceId := trace.TraceID{}
	spanId := trace.SpanID{}
	if traceparent != str.Empty {
		var err error
		s := strings.Split(traceparent, sym.Dash)
		traceId, err = trace.TraceIDFromHex(s[1])
		if err != nil {
			log.Warnf(parent, "Failed to get trace ID from traceparent, traceparent=%s, error=%v", traceparent, err)
		}
		spanId, err = trace.SpanIDFromHex(s[2])
		if err != nil {
			log.Warnf(parent, "Failed to get span ID from traceparent, traceparent=%s, error=%v", traceparent, err)
		}
	}

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
		SpanID:  spanId,
	})

	return trace.ContextWithSpanContext(parent, sc)
}

func TraceIdFromContext(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	return sc.TraceID().String()
}

func SpanIdFromContext(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	return sc.SpanID().String()
}

func ContextWithToken(parentContext context.Context, token string) context.Context {
	return context.WithValue(parentContext, tokenKey, token)
}

func TokenFromContext(ctx context.Context) (token string, exists bool) {
	if ctx.Value(tokenKey) == nil {
		return str.Empty, false
	}
	return ctx.Value(tokenKey).(string), true
}

func ContextWithUserId(parentContext context.Context, userId string) context.Context {
	return context.WithValue(parentContext, userIdKey, userId)
}

func UserIdFromContext(ctx context.Context) (userId string, exists bool) {
	if ctx.Value(userIdKey) == nil {
		return str.Empty, false
	}
	return ctx.Value(userIdKey).(string), true
}

func ContextWithUsername(parentContext context.Context, username string) context.Context {
	return context.WithValue(parentContext, usernameKey, username)
}

func UsernameFromContext(ctx context.Context) (username string, exists bool) {
	if ctx.Value(usernameKey) == nil {
		return str.Empty, false
	}
	return ctx.Value(usernameKey).(string), true
}

func ContextWithFullName(parentContext context.Context, fullname string) context.Context {
	return context.WithValue(parentContext, fullnameKey, fullname)
}

func FullNameFromContext(ctx context.Context) (fullname string, exists bool) {
	if ctx.Value(fullnameKey) == nil {
		return str.Empty, false
	}
	return ctx.Value(fullnameKey).(string), true
}

func ContextWithRealm(parentContext context.Context, realm string) context.Context {
	return context.WithValue(parentContext, realmKey, realm)
}

func RealmFromContext(ctx context.Context) (realm string, exists bool) {
	if ctx.Value(realmKey) == nil {
		return str.Empty, false
	}
	return ctx.Value(realmKey).(string), true
}

func ContextWithEmail(parentContext context.Context, email string) context.Context {
	return context.WithValue(parentContext, emailKey, email)
}

func EmailFromContext(ctx context.Context) (email string, exists bool) {
	if ctx.Value(emailKey) == nil {
		return str.Empty, false
	}
	return ctx.Value(emailKey).(string), true
}

func InjectMetadataToContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for _, key := range GetAllManagedKey() {
			values := md.Get(key)
			if len(values) > 0 && values[0] != str.Empty {
				ctx = context.WithValue(
					ctx,
					key,
					values[0],
				)
			}
		}
	}
	return ctx
}

func InjectContextToMetadata(ctx context.Context) context.Context {
	m := make(map[string]string)
	for _, key := range GetAllManagedKey() {
		m[key] = ctx.Value(key).(string)
	}
	md := metadata.New(m)
	return metadata.NewOutgoingContext(ctx, md)
}
