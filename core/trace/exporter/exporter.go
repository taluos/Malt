package exporter

import (
	otelExp "github.com/taluos/Malt/core/trace/exporter/otel"
	collectorExp "github.com/taluos/Malt/core/trace/exporter/otelCollector"
	zipkinExp "github.com/taluos/Malt/core/trace/exporter/zipkin"
	"github.com/taluos/Malt/pkg/errors"

	jaegerExp "github.com/taluos/Malt/core/trace/exporter/jaeger"

	fileExp "github.com/taluos/Malt/core/trace/exporter/file"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ExportConfig struct {
	Name     string  `json:",optional"`
	Endpoint string  `json:",optional"`
	Sampler  float64 `json:",default=1.0"`
	Batcher  string  `json:",default=jaeger"`
}

func CreateExporter(exportConfig ExportConfig, zipkinConfig []zipkin.Option, otelGrpcConfig []otlptracegrpc.Option, otelHttpConfig []otlptracehttp.Option) (trace.SpanExporter, error) {
	if len(exportConfig.Endpoint) > 0 {
		switch exportConfig.Batcher {
		case kindJaeger:
			return jaegerExp.NewJaegerExporter(exportConfig.Endpoint)
		case kindZipkin:
			return zipkinExp.NewZipkinExporter(exportConfig.Endpoint, zipkinConfig...)
		case kindotelgrpc:
			return otelExp.NewOtelGrpcExporter(exportConfig.Endpoint, otelGrpcConfig...)
		case kindotlphttp:
			return otelExp.NewOtelHttpExporter(exportConfig.Endpoint, otelHttpConfig...)
		case kindfile:
			return fileExp.NewFileExporter(exportConfig.Endpoint)
		case kindColletor:
			return collectorExp.Newotelcollector(exportConfig.Endpoint, otelHttpConfig...)
		default:
			return nil, errors.Errorf("unknow exporter type: %s", exportConfig.Batcher)
		}
	}

	return nil, errors.New("unsupport endpoint")

}
