// this is a modified version of jaeger exporter in  https://github.com/zeromicro/go-zero/blob/master/core/trace/agent.go
package jaeger

import (
	"Malt/pkg/errors"

	"net/url"

	"go.opentelemetry.io/otel/exporters/jaeger"
)

func NewJaegerExporter(endpoint string) (*jaeger.Exporter, error) {

	if endpoint == "" {
		return nil, errors.New("endpoint is empty")
	}

	u, err := url.Parse(endpoint)
	if err == nil && u.Scheme == "udp" {
		return jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(u.Hostname()),
			jaeger.WithAgentPort(u.Port())))
	}

	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))

}
