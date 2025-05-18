package main

import (
	"context"
	"time"

	"Malt/example/features/trace/rpc/service"
	pb "Malt/example/test_proto"
	"Malt/pkg/log"
	rpcserver "Malt/server/rpc/rpcServer"

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
	if err != nil {
		panic(err)
	}

	tp = traceSDK.NewTracerProvider(
		traceSDK.WithBatcher(jexp),
		traceSDK.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("test rpc server"),
				attribute.String("env", "test"),
			),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
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
	url := "http://192.168.142.140:14268/api/traces"
	NewTracerProvider(url)
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
