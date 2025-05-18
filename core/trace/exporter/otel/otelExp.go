package otel

import (
	"context"

	"github.com/taluos/Malt/pkg/errors"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func NewOtelHttpExporter(endpoint string, opts ...otlptracehttp.Option) (*otlptrace.Exporter, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint is empty")
	}

	opts = append(opts, otlptracehttp.WithEndpoint(endpoint))

	return otlptracehttp.New(context.Background(), opts...)
}

func NewOtelGrpcExporter(endpoint string, opts ...otlptracegrpc.Option) (*otlptrace.Exporter, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint is empty")
	}

	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))

	return otlptracegrpc.New(context.Background(), opts...)
}
