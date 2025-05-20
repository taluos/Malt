package middleware

import (
	maltAgent "github.com/taluos/Malt/core/trace"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(agent *maltAgent.Agent) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := maltAgent.NewTracer(trace.SpanKindServer,
			maltAgent.WithTracerProvider(agent.TracerProvider()),
			maltAgent.WithTracerName("http-handler"))

		carrier := propagation.HeaderCarrier(c.Request.Header)
		spanCtx, span := tr.Start(c.Request.Context(), c.FullPath(), agent.Propagator(), carrier)

		// 将span上下文传递给请求
		c.Request = c.Request.WithContext(spanCtx)

		// 处理请求
		c.Next()

		// 记录状态码和错误
		span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
		if len(c.Errors) > 0 {
			tr.End(spanCtx, span, c.Errors.Last().Err)
		} else {
			tr.End(spanCtx, span, nil)
		}
	}
}
