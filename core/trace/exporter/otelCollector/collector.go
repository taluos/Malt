package otelcollector

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func Newotelcollector(collectorURL string, otelHttpConfig ...otlptracehttp.Option) (*otlptrace.Exporter, error) {
	if collectorURL == "" {
		collectorURL = "http://localhost:4318" // default collector url
	}
	otelHttpConfig = append(otelHttpConfig, otlptracehttp.WithEndpointURL(collectorURL))
	return otlptracehttp.New(context.Background(), otelHttpConfig...)
}
