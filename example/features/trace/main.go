package main

import (
	"github.com/taluos/Malt/pkg/log"

	maltAgent "github.com/taluos/Malt/core/trace"
	rpcclient "github.com/taluos/Malt/example/features/trace/rpc/client"

	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

// 全局 TracerProvider
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// 创建全局 TracerProvider
	globalAgent = NewTracerProvider("trace-demo")
	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		err := globalAgent.Stop(ctx)
		if err != nil {
			panic(err)
		}
	}(ctx)

	tr := maltAgent.NewTracer(trace.SpanKindServer,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("func_main"),
	)

	spanCtx, span := tr.Start(ctx, "func_main", globalAgent.Propagator(), nil)
	defer tr.End(ctx, span, nil)

	var attr []attribute.KeyValue
	attr = append(attr, attribute.String("Key1", "value1"))

	span.SetAttributes(attr...)

	wg := &sync.WaitGroup{}

	wg.Add(5)

	go FuncA(spanCtx, wg)
	go FuncB(spanCtx, wg)
	go FuncC(spanCtx, wg)
	go FuncD(spanCtx, wg)
	go FuncE(spanCtx, wg)

	wg.Wait()

	span.AddEvent("this is a event")
	time.Sleep(time.Second * 2)
}

func FuncA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	// 使用全局 TracerProvider
	tr := maltAgent.NewTracer(trace.SpanKindServer,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("Internal"),
	)
	_, span := tr.Start(ctx, "Internal", globalAgent.Propagator(), nil)
	defer tr.End(ctx, span, err)

	span.SetAttributes(attribute.String("name", "Internal"))

	// do something

	type logStruct struct {
		User        string    `json:"user"`
		Auth        bool      `json:"auth"`
		CurrentTime time.Time `json:"currentTime"`
	}

	logTest := logStruct{
		User:        "user",
		Auth:        true,
		CurrentTime: time.Now(),
	}

	b, _ := json.Marshal(logTest)
	span.SetAttributes(attribute.Key("test log Key").String(string(b)))
	time.Sleep(time.Second * 1)
}

// http
func FuncB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8080/server")
	req.Header.SetMethod("GET")
	headers := make(map[string]string)

	var err error
	// 使用全局 TracerProvider
	tr := maltAgent.NewTracer(trace.SpanKindClient,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("Http"),
	)

	_, span := tr.Start(ctx, "Http", globalAgent.Propagator(), propagation.MapCarrier(headers))
	defer tr.End(ctx, span, err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	fcli := fasthttp.Client{}
	fesRes := fasthttp.Response{}
	_ = fcli.Do(req, &fesRes)

	time.Sleep(time.Second * 1)
}

// rpc
func FuncC(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	// 使用全局 TracerProvider
	tr := maltAgent.NewTracer(trace.SpanKindClient,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("Rpc"),
	)

	spanCtx, span := tr.Start(ctx, "Rpc", globalAgent.Propagator(), nil)
	defer tr.End(ctx, span, err)
	// 启动客户端

	if err = rpcclient.Run(spanCtx); err != nil {
		log.Fatalf("客户端运行失败: %v", err)
	}

	time.Sleep(2 * time.Second)
}

// gorm
func FuncD(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8080/gorm")
	req.Header.SetMethod("GET")

	// 创建 headers 映射用于存储追踪信息
	headers := make(map[string]string)

	// 使用全局 TracerProvider
	tr := maltAgent.NewTracer(trace.SpanKindClient,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("Gorm"),
	)

	// 使用 MapCarrier 传递追踪上下文
	_, span := tr.Start(ctx, "Gorm", globalAgent.Propagator(), propagation.MapCarrier(headers))
	defer tr.End(ctx, span, nil)

	// 将追踪信息设置到请求头中
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	fcli := fasthttp.Client{}
	fesRes := fasthttp.Response{}
	_ = fcli.Do(req, &fesRes)

	time.Sleep(time.Second * 1)
}

// redis
func FuncE(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	// 使用全局 TracerProvider
	tr := maltAgent.NewTracer(trace.SpanKindClient,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("Redis"),
	)

	spanCtx, span := tr.Start(ctx, "Redis", globalAgent.Propagator(), nil)
	defer tr.End(ctx, span, err)

	redisDB := redis.NewClient(&redis.Options{
		Addr: "192.168.142.137:6379",
	})

	// Enable tracing instrumentation.
	if err = redisotel.InstrumentTracing(redisDB); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err = redisotel.InstrumentMetrics(redisDB); err != nil {
		panic(err)
	}

	redisDB.Set(spanCtx, "name", "Redis otel test", time.Minute*5)

	time.Sleep(time.Second * 1)

}
