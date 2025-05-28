package clientinterceptors

import (
	"context"

	maltAgent "github.com/taluos/Malt/core/trace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryTracingInterceptor(agent *maltAgent.Agent) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tr := maltAgent.NewTracer(trace.SpanKindClient,
			maltAgent.WithTracerProvider(agent.TracerProvider()),
			maltAgent.WithTracerName(method))
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		carrier := propagation.HeaderCarrier(md)
		spanCtx, span := tr.Start(ctx, method, agent.Propagator(), carrier)
		err := invoker(spanCtx, method, req, reply, cc, opts...)
		tr.End(spanCtx, span, err)
		return err
	}
}

func StreamTracingInterceptor(agent *maltAgent.Agent) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		tr := maltAgent.NewTracer(trace.SpanKindClient,
			maltAgent.WithTracerProvider(agent.TracerProvider()),
			maltAgent.WithTracerName(method))
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		carrier := propagation.HeaderCarrier(md)
		spanCtx, span := tr.Start(ctx, method, agent.Propagator(), carrier)
		clientStream, err := streamer(spanCtx, desc, cc, method, opts...)
		tr.End(spanCtx, span, err)
		return clientStream, err
	}
}
