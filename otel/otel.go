package otel

import (
	"context"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	myLog "github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.14.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	errOtelInitIsDisabled = "Otel init is disabled"
)

var (
	cfg      *config.Object
	tracer   trace.Tracer
	meter    metric.Meter
	counters map[string]metric.Int64Counter
)

// Init Initializes an OTLP exporter, and configures the corresponding trace and metric providers.
func Init(ctx context.Context, config config.Config) {
	cfg = config.GetObject()
	if isOtelConfigMissingOrDisabled(cfg) {
		log.Info(ctx, errOtelInitIsDisabled)
		return
	}
	serviceName := cfg.App.Name
	var err error
	var spanExporter sdktrace.SpanExporter
	var metricExporter sdkmetric.Exporter

	// init otel trace exporter
	traceConfig := cfg.Otel.Trace
	if traceConfig.Exporter.Disabled {
		log.Warn(ctx, "Otel trace exporter is disabled")
	} else {
		exporterType := traceConfig.Exporter.Type
		switch exporterType {
		case exporterTypeOtlpGrpc:
			// set up a trace OTLP gRPC exporter
			exporterConfig := traceConfig.Exporter.Otlp.Grpc
			conn, cancel, err := initGrpcConn(ctx, exporterConfig)
			handleErr(ctx, err, "Failed to create gRPC connection for trace OTLP exporter on otel init")
			defer cancel()

			spanExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
			handleErr(ctx, err, "Failed to create collector trace exporter on otel init")

		default:
			log.Fatalf(ctx, response.UnsupportedType, "Unsupported trace exporter type, type=%s", exporterType)
			return
		}
	}

	// init otel metric exporter
	metricConfig := cfg.Otel.Metric
	if metricConfig.Exporter.Disabled {
		log.Warn(ctx, "Otel metric exporter is disabled")
	} else {
		exporterType := metricConfig.Exporter.Type
		switch exporterType {
		case exporterTypeOtlpGrpc:
			// Set up a metric OTLP gRPC exporter
			exporterConfig := metricConfig.Exporter.Otlp.Grpc
			conn, cancel, err := initGrpcConn(ctx, exporterConfig)
			handleErr(ctx, err, "Failed to create gRPC connection for metric OTLP exporter on otel init")
			defer cancel()

			metricExporter, err = otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
			handleErr(ctx, err, "Failed to create collector metric exporter on otel init")

		default:
			log.Fatalf(ctx, response.UnsupportedType, "Unsupported metric exporter type, type=%s", exporterType)
			return
		}
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(cfg.App.Name),
		),
	)
	handleErr(ctx, err, "failed to create resource on otel init")

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))

	// set tracer provider
	if spanExporter != nil {
		bsp := sdktrace.NewBatchSpanProcessor(spanExporter)
		tracerProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
		)
		otel.SetTracerProvider(tracerProvider)
	}

	// set metric provider
	if metricExporter != nil {
		meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)))
		otel.SetMeterProvider(meterProvider)
	}

	// Set default tracer
	tracer = otel.Tracer(serviceName + "-tracer")

	// Set default meter
	meter = otel.Meter(cfg.Otel.Metric.InstrumentationName)

	// Init default counters
	counters = make(map[string]metric.Int64Counter)
}

func initGrpcConn(ctx context.Context, exporterConfig *config.OtelExporterOtlpGrpcConfig) (*grpc.ClientConn, context.CancelFunc, error) {
	opts := make([]grpc.DialOption, integer.Zero)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	clientMaxReceiveMessageSizeStr := exporterConfig.ClientMaxReceiveMessageSize
	if clientMaxReceiveMessageSizeStr != str.Empty {
		clientMaxReceiveMessageSize, err := bytesize.Parse(clientMaxReceiveMessageSizeStr)
		if err != nil {
			log.Fatalf(ctx, err, "Failed to parse otel trace OTLP gRPC client max received message size '%s', %s", clientMaxReceiveMessageSizeStr, err.Error())
			return nil, nil, err
		}
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(clientMaxReceiveMessageSize))))
	}

	// Create GRPC connection with timeout
	ctxCancel, cancel := context.WithTimeout(ctx, time.Duration(exporterConfig.Timeout)*time.Second)
	conn, err := grpc.DialContext(
		ctxCancel,
		exporterConfig.Address,
		opts...,
	)
	return conn, cancel, err
}

func handleErr(ctx context.Context, err error, message string) {
	if err != nil {
		log.Fatal(ctx, err, message)
		return
	}
}

func Trace(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, *SpanWrapper) {
	// Get parent span if any
	sc := trace.SpanContextFromContext(ctx)
	ctx = context.WithValue(ctx, myLog.SpanParentIdLogKey, sc.SpanID().String())

	var span trace.Span
	if isOtelConfigMissingOrDisabled(cfg) {
		return ctx, &SpanWrapper{}
	} else {
		ctx, span = tracer.Start(ctx, spanName, opts...)
	}
	return ctx, &SpanWrapper{span}
}

// End completes the Span. The Span is considered complete and ready to be
// delivered through the rest of the telemetry pipeline after this method
// is called. Therefore, updates to the Span are not allowed after this
// method has been called.
func (w *SpanWrapper) End(options ...trace.SpanEndOption) {
	if !isOtelConfigMissingOrDisabled(cfg) {
		w.span.End(options...)
	}
}

// AddEvent adds an event with the provided name and options.
func (w *SpanWrapper) AddEvent(name string, options ...trace.EventOption) {
	if !isOtelConfigMissingOrDisabled(cfg) {
		w.span.AddEvent(name, options...)
	}
}

// IsRecording returns the recording state of the Span. It will return
// true if the Span is active and events can be recorded.
func (w *SpanWrapper) IsRecording() bool {
	if !isOtelConfigMissingOrDisabled(cfg) {
		return w.span.IsRecording()
	}
	return false
}

// RecordError will record err as an exception span event for this span. An
// additional call to SetStatus is required if the Status of the Span should
// be set to Error, as this method does not change the Span status. If this
// span is not being recorded or err is nil then this method does nothing.
func (w *SpanWrapper) RecordError(err error, options ...trace.EventOption) {
	if !isOtelConfigMissingOrDisabled(cfg) {
		w.span.RecordError(err, options...)
	}
}

// SpanContext returns the SpanContext of the Span. The returned SpanContext
// is usable even after the End method has been called for the Span.
func (w *SpanWrapper) SpanContext() trace.SpanContext {
	if !isOtelConfigMissingOrDisabled(cfg) {
		return w.span.SpanContext()
	}
	return trace.SpanContext{}
}

// SetStatus sets the status of the Span in the form of a code and a
// description, overriding previous values set. The description is only
// included in a status when the code is for an error.
func (w *SpanWrapper) SetStatus(code codes.Code, description string) {
	if !isOtelConfigMissingOrDisabled(cfg) {
		w.span.SetStatus(code, description)
	}
}

// SetName sets the Span name.
func (w *SpanWrapper) SetName(name string) {
	if !isOtelConfigMissingOrDisabled(cfg) {
		w.span.SetName(name)
	}
}

// SetAttributes sets kv as attributes of the Span. If a key from kv
// already exists for an attribute of the Span it will be overwritten with
// the value contained in kv.
func (w *SpanWrapper) SetAttributes(kv ...attribute.KeyValue) {
	if !isOtelConfigMissingOrDisabled(cfg) {
		w.span.SetAttributes(kv...)
	}
}

// TracerProvider returns a TracerProvider that can be used to generate
// additional Spans on the same telemetry pipeline as the current Span.
func (w *SpanWrapper) TracerProvider() trace.TracerProvider {
	if !isOtelConfigMissingOrDisabled(cfg) {
		return w.span.TracerProvider()
	}
	return nil
}

func AddCounter(_ context.Context, counterName string, unit string) error {
	counter, err := meter.Int64Counter(counterName, metric.WithUnit(unit))
	if err != nil {
		return err
	}

	counters[counterName] = counter
	return nil
}

func Count(ctx context.Context, counterName string, incr int64, opts ...metric.AddOption) {
	counters[counterName].Add(ctx, incr, opts...)
}

func isOtelConfigMissingOrDisabled(config *config.Object) bool {
	return config == nil || config.Otel == nil || config.Otel.Disabled
}
