package middleware

import (
	"net/http"

	maltAgent "github.com/taluos/Malt/core/trace"
	"github.com/valyala/fasthttp"

	fiber "github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// FastHTTPHeadersToHTTPHeaders 将fasthttp.RequestHeader转换为http.Header
func FastHTTPHeadersToHTTPHeaders(fh *fasthttp.Request) http.Header {
	h := make(http.Header)
	fh.Header.VisitAll(func(key, value []byte) {
		h.Set(string(key), string(value))
	})
	return h
}

func TracingMiddleware(agent *maltAgent.Agent) fiber.Handler {
	return func(c fiber.Ctx) error {
		tr := maltAgent.NewTracer(trace.SpanKindServer,
			maltAgent.WithTracerProvider(agent.TracerProvider()),
			maltAgent.WithTracerName(c.Path()))

		// 将fasthttp.RequestHeader转换为http.Header
		httpHeaders := FastHTTPHeadersToHTTPHeaders(c.Request())
		carrier := propagation.HeaderCarrier(httpHeaders)

		spanCtx, span := tr.Start(c.Context(),
			c.Path(),
			agent.Propagator(),
			carrier)

		// 将span上下文传递给请求
		c.SetContext(spanCtx)

		// 处理请求
		err := c.Next()

		// 记录状态码和错误
		span.SetAttributes(attribute.Int("http.status_code", c.Response().StatusCode()))
		if err != nil {
			tr.End(spanCtx, span, err)
		} else {
			tr.End(spanCtx, span, nil)
		}
		return err
	}
}
