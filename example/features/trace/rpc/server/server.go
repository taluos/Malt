package main

import (
	"context"
	"time"

	maltAgent "github.com/taluos/Malt/core/trace"
	"github.com/taluos/Malt/example/features/trace/rpc/service"
	pb "github.com/taluos/Malt/example/test_proto"
	"github.com/taluos/Malt/pkg/log"
	rpcserver "github.com/taluos/Malt/server/rpc/rpcServer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func NewTracerProvider(name string) *maltAgent.Agent {

	agent := maltAgent.NewAgent(name, "http://localhost:4318", "ratio", 1.0, "collector",
		maltAgent.WithTracerProviderOptions(traceSDK.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			attribute.String("env", "test"),
		))),
	)

	return agent
}

func rpcServerInit() *rpcserver.Server {
	// 创建 gRPC Server，可根据需要自定义监听地址、超时时间等
	s := rpcserver.NewServer(
		rpcserver.WithAddress("127.0.0.1:50051"),
		rpcserver.WithTimeout(5*time.Second),
		rpcserver.WithEnableTracing(true),
	)

	// 注册服务
	pb.RegisterGreeterServer(s.Server, &service.GreeterServer{})

	return s
}

func rpcRun(srv *rpcserver.Server, ctx context.Context) error {
	return srv.Start(ctx)
}

func rpcStop(srv *rpcserver.Server, ctx context.Context) error {
	<-ctx.Done() // 等待接收退出信号
	log.Info("RPC server stopping")
	sctx, cancal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancal()
	return srv.Stop(sctx)
}

func main() {
	var err error
	agent := NewTracerProvider("Rpc Server")
	defer agent.Stop(context.Background())
	tr := maltAgent.NewTracer(trace.SpanKindServer,
		maltAgent.WithTracerProvider(agent.TracerProvider()),
		maltAgent.WithTracerName("test server"),
	)

	spanCtx, span := tr.Start(context.Background(), "test", agent.Propagator(), nil)
	defer tr.End(context.Background(), span, err)

	rpcServer := rpcServerInit()

	log.Info("RPC server starting")
	err = rpcRun(rpcServer, spanCtx)
	if err != nil {
		log.Fatalf("RPC server stopped with error: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = rpcStop(rpcServer, spanCtx)
	if err != nil {
		panic(err)
	}

}
