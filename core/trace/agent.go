package trace

import (
	"sync"

	"github.com/taluos/Malt/pkg/errors"

	"github.com/taluos/Malt/core/trace/exporter"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type Agent struct {
	Name        string  `json:",optional"`
	Endpoint    string  `json:",optional"`
	SamplerMode string  `json:",optional,default=never,options=ratio|always|never"`
	Sampler     float64 `json:",default=1.0"`
	Batcher     string  `json:",default=jaeger,options=zipkin|jaeger|prometheus|otelgrpc|otelhttp|file"`

	opt *telemetryOptions
}

var (
	lock sync.Mutex
)

func NewAgent(name string, endpoint string, samplerMode string, sampler float64, batcher string, o ...TelemetryOptions) *Agent {
	opt := &telemetryOptions{}

	for _, f := range o {
		f(opt)
	}

	return &Agent{
		Name:        name,
		Endpoint:    endpoint,
		SamplerMode: samplerMode,
		Sampler:     sampler,
		Batcher:     batcher,

		opt: opt,
	}
}

func InitAgent(agent *Agent) *trace.TracerProvider {
	lock.Lock()
	defer lock.Unlock()

	tp, err := agent.startAgent()
	if err != nil {
		panic(err)
	}
	return tp
}

func (agent *Agent) startAgent() (*trace.TracerProvider, error) {
	var (
		exp trace.SpanExporter
		err error

		tracerOptions []trace.TracerProviderOption
	)

	exporterConfig := exporter.ExportConfig{
		Name:     agent.Name,
		Endpoint: agent.Endpoint,
		Sampler:  agent.Sampler,
		Batcher:  agent.Batcher,
	}

	exp, err = exporter.CreateExporter(exporterConfig,
		agent.opt.zipkinOptions,
		agent.opt.otelGrpcOptions,
		agent.opt.otelHttpOptions,
	)
	if err != nil {
		return nil, err
	}

	tracerOptions = append(tracerOptions,
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(agent.Name))),
	)

	switch agent.SamplerMode {
	case "ratio":
		tracerOptions = append(tracerOptions, trace.WithSampler(trace.TraceIDRatioBased(agent.Sampler)))
	case "always":
		tracerOptions = append(tracerOptions, trace.WithSampler(trace.AlwaysSample()))
	case "never":
		tracerOptions = append(tracerOptions, trace.WithSampler(trace.NeverSample()))
	default:
		return nil, errors.Errorf("invalid sampler mode: %s", agent.SamplerMode)
	}

	tracerOptions = append(tracerOptions, agent.opt.tracerProviderOptions...)

	tp := trace.NewTracerProvider(tracerOptions...)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		errors.Errorf("[otel] error: %v", err)
	}))

	return tp, nil
}
