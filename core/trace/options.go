package trace

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

// telemetryOptions is an opentelelmetry configs.
type telemetryOptions struct {
	// exporter options
	zipkinOptions    []zipkin.Option
	otelGrpcOptions  []otlptracegrpc.Option
	otelHttpOptions  []otlptracehttp.Option
	collectorOptions []otlptracehttp.Option

	// tracer provider options
	tracerProviderOptions []sdkTrace.TracerProviderOption
}

type tracerOptions struct {
	tp *sdkTrace.TracerProvider

	Name       string
	Attributes []attribute.KeyValue
	Resource   *resource.Resource
}

type TelemetryOptions func(*telemetryOptions)
type TracerOptions func(*tracerOptions)

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

func WithTracerProviderOptions(opts ...sdkTrace.TracerProviderOption) TelemetryOptions {
	return func(o *telemetryOptions) {
		o.tracerProviderOptions = opts
	}
}

func WithTracerProvider(tp *sdkTrace.TracerProvider) TracerOptions {
	return func(o *tracerOptions) {
		o.tp = tp
	}
}

func WithTracerName(name string) TracerOptions {
	return func(o *tracerOptions) {
		o.Name = name
	}
}

func WithTracerAttributes(attrs ...attribute.KeyValue) TracerOptions {
	return func(o *tracerOptions) {
		o.Attributes = attrs
	}
}

func WithTracerResource(res *resource.Resource) TracerOptions {
	return func(o *tracerOptions) {
		o.Resource = res
	}
}
