package client

import (
	pb "Malt/example/test_proto"
	rpcclient "Malt/server/rpc/rpcClient"

	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	tp *traceSDK.TracerProvider
)

func NewTracerProvider(url string) {
	collectorURL := "http://localhost:4318" // Collector 的默认 OTLP HTTP 端点
	jexp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpointURL(collectorURL),
	)
	//jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	// fmt.Println(jexp)
	if err != nil {
		panic(err)
	}

	tp = traceSDK.NewTracerProvider(
		traceSDK.WithBatcher(jexp),
		traceSDK.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("test rpc client"),
				attribute.String("env", "test"),
			),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

// Run 启动 gRPC 客户端并优雅关闭
func Run(ctx context.Context) error {
	url := "http://192.168.142.140:14268/api/traces"

	NewTracerProvider(url)
	tr := tp.Tracer("test")

	spanCtx, span := tr.Start(ctx, "test")

	time.Sleep(time.Second * 1)
	defer span.End()
	// 创建 gRPC 客户端，可根据需要自定义连接地址、超时时间等
	c, err := rpcclient.NewClient(
		rpcclient.WithEndpoint("127.0.0.1:50051"),
		rpcclient.WithTimeout(5*time.Second),
		rpcclient.WithInsecure(true),
		rpcclient.WithEnableTracing(true),
		// rpcclient.WithUnaryInterceptors(otelgrpc.UnaryClientInterceptor()),
		// rpcclient.WithOptions(grpc.WithStatsHandler(otelgrpc.NewClientHandler())),
	)

	if err != nil {
		return err
	}

	// 创建 Greeter 客户端
	client := pb.NewGreeterClient(c.CliConn)

	// 调用 SayHello 方法
	resp, err := client.SayHello(spanCtx, &pb.HelloRequest{Name: "Malt用户"})
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
