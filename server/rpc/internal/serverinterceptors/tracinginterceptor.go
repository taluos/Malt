package serverinterceptors

import (
	"context"

	maltAgent "github.com/taluos/Malt/core/trace"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TracingUnaryServerInterceptor(agent *maltAgent.Agent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tr := maltAgent.NewTracer(trace.SpanKindServer,
			maltAgent.WithTracerProvider(agent.TracerProvider()),
			maltAgent.WithTracerName("rpc-handler"))

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		carrier := propagation.HeaderCarrier(md)

		spanCtx, span := tr.Start(ctx, info.FullMethod, agent.Propagator(), carrier)

		resp, err := handler(spanCtx, req)

		tr.End(spanCtx, span, err)

		return resp, err
	}
}
