package main

import (
	"github.com/taluos/Malt/pkg/log"

	rpcclient "github.com/taluos/Malt/example/features/trace/rpc/client"

	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"golang.org/x/net/context"
)

var (
	Tp *trace.TracerProvider
)

func NewTracerProvider() (*trace.TracerProvider, error) { //url string
	collectorURL := "http://localhost:4318" // Collector 的默认 OTLP HTTP 端点
	jexp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpointURL(collectorURL),
	)

	// jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	Tp = trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("test"),
				attribute.String("env", "test"),
			),
		),
	)
	otel.SetTracerProvider(Tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return Tp, nil
}

func main() {
	// url := "http://192.168.142.140:14268/api/traces"
	ctx, cancel := context.WithCancel(context.Background())

	tp, _ := NewTracerProvider()

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		err := tp.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}(ctx)

	tr := tp.Tracer("test")
	defer tp.Shutdown(context.Background())
	spanCtx, span := tr.Start(ctx, "func_main")

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
	defer span.End()
}

func FuncA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	tr := Tp.Tracer("Func A")
	_, span := tr.Start(ctx, "Func A !!!!")
	defer span.End()

	span.SetAttributes(attribute.String("name", "func_A"))

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

	tr := Tp.Tracer("Func B")

	spanCtx, span := tr.Start(ctx, "Func B!!!!")

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8080/server")
	req.Header.SetMethod("GET")

	// 传播器
	p := otel.GetTextMapPropagator()

	// 包裹
	headers := make(map[string]string)
	p.Inject(spanCtx, propagation.MapCarrier(headers))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	fcli := fasthttp.Client{}
	fesRes := fasthttp.Response{}
	_ = fcli.Do(req, &fesRes)

	time.Sleep(time.Second * 1)
	defer span.End()
}

// rpc
func FuncC(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer("Func C")

	spanCtx, span := tr.Start(ctx, "Func C!!!!")

	// 启动客户端

	if err := rpcclient.Run(spanCtx); err != nil {
		log.Fatalf("客户端运行失败: %v", err)
	}

	time.Sleep(2 * time.Second)

	defer span.End()
}

// gorm
func FuncD(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	tr := Tp.Tracer("Func D")

	spanCtx, span := tr.Start(ctx, "Func D!!!!")

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8080/gorm")
	req.Header.SetMethod("GET")

	// 传播器
	p := otel.GetTextMapPropagator()

	// 包裹
	headers := make(map[string]string)
	p.Inject(spanCtx, propagation.MapCarrier(headers))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	fcli := fasthttp.Client{}
	fesRes := fasthttp.Response{}
	_ = fcli.Do(req, &fesRes)

	time.Sleep(time.Second * 1)
	defer span.End()
}

// redis
func FuncE(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	tr := otel.Tracer("Func E")

	spanCtx, span := tr.Start(ctx, "Func E!!!!")

	redisDB := redis.NewClient(&redis.Options{
		Addr: "192.168.142.137:6379",
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(redisDB); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(redisDB); err != nil {
		panic(err)
	}

	defer span.End()

	redisDB.Set(spanCtx, "name", "Redis otel test", time.Minute*5)

	time.Sleep(time.Second * 1)

}
