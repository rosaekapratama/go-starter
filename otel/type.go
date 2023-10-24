package otel

import "go.opentelemetry.io/otel/trace"

type SpanWrapper struct {
	span trace.Span
}
