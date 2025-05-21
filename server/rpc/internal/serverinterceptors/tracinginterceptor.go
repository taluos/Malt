package serverinterceptors

import (
	"context"

	maltAgent "github.com/taluos/Malt/core/trace"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryTracingInterceptor(agent *maltAgent.Agent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		tr := maltAgent.NewTracer(trace.SpanKindServer,
			maltAgent.WithTracerProvider(agent.TracerProvider()),
			maltAgent.WithTracerName(info.FullMethod))

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		carrier := propagation.HeaderCarrier(md)

		spanCtx, span := tr.Start(ctx, info.FullMethod, agent.Propagator(), carrier)

		resp, err := handler(spanCtx, req)

		defer tr.End(spanCtx, span, err)

		return resp, err
	}
}
