package client

import (
	"context"
	"log"
	"time"

	maltAgent "github.com/taluos/Malt/core/trace"
	pb "github.com/taluos/Malt/example/test_proto"
	rpcclient "github.com/taluos/Malt/server/rpc/rpcClient"

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

// Run 启动 gRPC 客户端并优雅关闭
func Run(ctx context.Context) error {
	var err error
	globalAgent := NewTracerProvider("Rpc Client")
	defer globalAgent.Shutdown(ctx)

	time.Sleep(time.Second * 1)
	// 创建 gRPC 客户端，可根据需要自定义连接地址、超时时间等
	c, err := rpcclient.NewClient(
		rpcclient.WithEndpoint("127.0.0.1:50051"),
		rpcclient.WithTimeout(5*time.Second),
		rpcclient.WithInsecure(true),
		rpcclient.WithEnableTracing(true),
		rpcclient.WithAgent(globalAgent),
	)

	if err != nil {
		return err
	}

	// 创建 Greeter 客户端
	client := pb.NewGreeterClient(c.ClientConn)

	// 调用 SayHello 方法
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Malt用户"})
	if err != nil {
		log.Printf("调用 SayHello 失败: %v", err)
	} else {
		log.Printf("收到服务器响应: %s", resp.GetMessage())
	}

	if err := c.Close(context.Background()); err != nil {
		return err
	}

	return nil
}
