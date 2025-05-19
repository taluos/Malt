package trace

import (
	"context"

	"github.com/taluos/Malt/pkg/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	kind   trace.SpanKind
	tracer trace.Tracer

	opts *tracerOptions
}

func NewTracer(kind trace.SpanKind, opts ...TracerOptions) *Tracer {

	o := tracerOptions{
		Name:       defaultTracerName,
		Attributes: []attribute.KeyValue{},
		Resource:   resource.Default(),
	}

	for _, opt := range opts {
		opt(&o)
	}

	switch kind {
	case trace.SpanKindClient:
		return &Tracer{
			kind:   kind,
			tracer: o.tp.Tracer(o.Name),
			opts:   &o,
		}
	case trace.SpanKindServer:
		return &Tracer{
			kind:   kind,
			tracer: o.tp.Tracer(o.Name),
			opts:   &o,
		}
	default:
		log.Warnf("unknown span kind: %s", kind)
		return &Tracer{
			kind:   trace.SpanKindInternal,
			tracer: o.tp.Tracer(o.Name),
			opts:   &o,
		}
	}
}

func (t *Tracer) Start(ctx context.Context, spanName string, propagator propagation.TextMapPropagator, carrier propagation.TextMapCarrier) (context.Context, trace.Span) {
	if carrier == nil {
		carrier = propagation.MapCarrier{}
	}
	if t.kind == trace.SpanKindServer {
		ctx = propagator.Extract(ctx, carrier)
	}
	spanCtx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(t.kind))

	if t.kind == trace.SpanKindClient {
		propagator.Inject(spanCtx, carrier)
	}

	return spanCtx, span
}

func (t *Tracer) End(ctx context.Context, span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.SetStatus(codes.Ok, "OK")

	span.End()
}
