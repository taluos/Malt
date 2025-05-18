package trace

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/trace"
)

// telemetryOptions is an opentelelmetry configs.
type telemetryOptions struct {
	// exporter options
	zipkinOptions    []zipkin.Option
	otelGrpcOptions  []otlptracegrpc.Option
	otelHttpOptions  []otlptracehttp.Option
	collectorOptions []otlptracehttp.Option

	// tracer provider options
	tracerProviderOptions []trace.TracerProviderOption
}

type TelemetryOptions func(*telemetryOptions)

func WithZipkinOptions(opts ...zipkin.Option) TelemetryOptions {
	return func(o *telemetryOptions) {
		o.zipkinOptions = opts
	}
}

func WithOtelGrpcOptions(opts ...otlptracegrpc.Option) TelemetryOptions {
	return func(o *telemetryOptions) {
		o.otelGrpcOptions = opts
	}
}

func WithOtelHttpOptions(opts ...otlptracehttp.Option) TelemetryOptions {
	return func(o *telemetryOptions) {
		o.otelHttpOptions = opts
	}
}

func WithCollectorOptions(opts ...otlptracehttp.Option) TelemetryOptions {
	return func(o *telemetryOptions) {
		o.collectorOptions = opts
	}
}

func WithTracerProviderOptions(opts ...trace.TracerProviderOption) TelemetryOptions {
	return func(o *telemetryOptions) {
		o.tracerProviderOptions = opts
	}
}
