package zipkin

import (
	"github.com/taluos/Malt/pkg/errors"

	"go.opentelemetry.io/otel/exporters/zipkin"
)

func NewZipkinExporter(endpoint string, opts ...zipkin.Option) (*zipkin.Exporter, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint is empty")
	}
	return zipkin.New(endpoint, opts...)
}
