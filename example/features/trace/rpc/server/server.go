package main

import (
	"context"
	"time"

	maltAgent "github.com/taluos/Malt/core/trace"
	"github.com/taluos/Malt/example/features/trace/rpc/service"
	pb "github.com/taluos/Malt/example/test_proto"
	"github.com/taluos/Malt/pkg/log"
	rpcserver "github.com/taluos/Malt/server/rpc"
	grpcServer "github.com/taluos/Malt/server/rpc/rpc-grpc"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var globalAgent *maltAgent.Agent

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

func rpcRun(srv rpcserver.Server, ctx context.Context) error {
	return srv.Start(ctx)
}

func rpcStop(srv rpcserver.Server, ctx context.Context) error {
	<-ctx.Done() // 等待接收退出信号
	log.Info("RPC server stopping")
	sctx, cancal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancal()
	return srv.Stop(sctx)
}

func main() {
	var err error
	ctx := context.Background()
	globalAgent := NewTracerProvider("Rpc Server")
	defer globalAgent.Shutdown(context.Background())

	rpcServer := rpcserver.NewServer("grpc",
		grpcServer.WithServerAddress("127.0.0.1:50051"),
		grpcServer.WithServerTimeout(5*time.Second),
		grpcServer.WithServerEnableTracing(true),
		grpcServer.WithServerAgent(globalAgent),
	)

	rpcServer = rpcServer.RegisterService(pb.RegisterGreeterServer, &service.GreeterServer{})
	// 注册服务
	// pb.RegisterGreeterServer(rpcServer.Server, &service.GreeterServer{})

	log.Info("RPC server starting")
	err = rpcRun(rpcServer, ctx)
	if err != nil {
		log.Fatalf("RPC server stopped with error: %v", err)
	}

	err = rpcStop(rpcServer, ctx)
	if err != nil {
		panic(err)
	}

}
