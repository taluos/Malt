package trace

import (
	"context"
	"sync"

	"github.com/taluos/Malt/pkg/errors"

	"github.com/taluos/Malt/core/trace/exporter"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type Agent struct {
	Name        string  `json:",optional"`
	Endpoint    string  `json:",optional"`
	SamplerMode string  `json:",optional,default=never,options=ratio|always|never"`
	Sampler     float64 `json:",default=1.0"`
	Batcher     string  `json:",default=jaeger,options=zipkin|jaeger|prometheus|otelgrpc|otelhttp|file"`

	tp         *sdkTrace.TracerProvider
	propagator propagation.TextMapPropagator
	errHandler otel.ErrorHandler
	mu         sync.Mutex

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

	agent := &Agent{
		Name:        name,
		Endpoint:    endpoint,
		SamplerMode: samplerMode,
		Sampler:     sampler,
		Batcher:     batcher,

		opt: opt,
	}

	err := agent.startAgent()
	if err != nil {
		panic(err)
	}

	return agent

}

func (agent *Agent) startAgent() error {
	var (
		exp sdkTrace.SpanExporter
		err error

		tracerOptions []sdkTrace.TracerProviderOption
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
		return err
	}

	tracerOptions = append(tracerOptions,
		sdkTrace.WithBatcher(exp),
		sdkTrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(agent.Name))),
	)

	switch agent.SamplerMode {
	case "ratio":
		tracerOptions = append(tracerOptions, sdkTrace.WithSampler(sdkTrace.TraceIDRatioBased(agent.Sampler)))
	case "always":
		tracerOptions = append(tracerOptions, sdkTrace.WithSampler(sdkTrace.AlwaysSample()))
	case "never":
		tracerOptions = append(tracerOptions, sdkTrace.WithSampler(sdkTrace.NeverSample()))
	default:
		return errors.Errorf("invalid sampler mode: %s", agent.SamplerMode)
	}

	tracerOptions = append(tracerOptions, agent.opt.tracerProviderOptions...)

	tp := sdkTrace.NewTracerProvider(tracerOptions...)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	errHandler := otel.ErrorHandlerFunc(func(err error) {
		errors.Errorf("[otel] error: %v", err)
	})

	agent.mu.Lock()
	agent.tp = tp
	agent.propagator = propagator
	agent.errHandler = errHandler
	agent.mu.Unlock()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(agent.propagator)
	otel.SetErrorHandler(agent.errHandler)

	return nil
}

func (agent *Agent) Stop(ctx context.Context) error {
	agent.mu.Lock()
	defer agent.mu.Unlock()

	if agent.tp != nil {
		err := agent.tp.Shutdown(ctx)
		agent.tp = nil
		return err
	}
	return nil
}

func (agent *Agent) TracerProvider() *sdkTrace.TracerProvider {
	return agent.tp
}

func (agent *Agent) Propagator() propagation.TextMapPropagator {
	return agent.propagator
}

func (agent *Agent) ErrorHandler() otel.ErrorHandler {
	return agent.errHandler
}
