package main

import (
	"context"
	"time"

	agent "github.com/taluos/Malt/core/trace"
	"github.com/taluos/Malt/example/features/trace/rpc/service"
	pb "github.com/taluos/Malt/example/test_proto"
	"github.com/taluos/Malt/pkg/log"
	rpcserver "github.com/taluos/Malt/server/rpc/rpcServer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func NewTracerProvider(name string) *traceSDK.TracerProvider {

	agentOpt := agent.NewAgent(name, "http://localhost:4318", "ratio", 1.0, "collector",
		agent.WithTracerProviderOptions(traceSDK.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			attribute.String("env", "test"),
		))),
	)

	tp := agent.InitAgent(agentOpt)
	return tp
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

	tp := NewTracerProvider("Rpc Server")
	defer tp.Shutdown(context.Background())
	tr := tp.Tracer("test")

	spanCtx, span := tr.Start(context.Background(), "test")
	defer span.End()

	rpcServer := rpcServerInit()

	log.Info("RPC server starting")
	err := rpcRun(rpcServer, spanCtx)
	if err != nil {
		log.Fatalf("RPC server stopped with error: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = rpcStop(rpcServer, spanCtx)
	if err != nil {
		panic(err)
	}

}
