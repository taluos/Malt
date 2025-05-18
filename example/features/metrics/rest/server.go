package main

import (
	"context"
	"log"

	httpserver "github.com/taluos/Malt/server/rest/httpServer"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {

	gaugeMetric := &ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        "example_gauge_metric",
		Description: "an example of gauge type metric",
		Labels:      []string{"label1"},
	}

	// 在服务器启动前注册指标
	monitor := ginmetrics.GetMonitor()
	if err := monitor.AddMetric(gaugeMetric); err != nil {
		log.Fatalf("添加指标失败: %v", err)
	}

	r := httpserver.NewServer(
		httpserver.WithPort(8090),
		httpserver.WithEnableProfiling(true),
		httpserver.WithEnableTracing(true),
		httpserver.WithEnableMetrics(true),
		httpserver.WithMiddleware(gin.Recovery()),
	)

	r.GET("/test/:id", func(ctx *gin.Context) {
		// 更新指标值
		_ = ginmetrics.GetMonitor().GetMetric("example_gauge_metric").Inc([]string{"label_value1"})

		ctx.JSON(200, map[string]string{
			"productId": ctx.Param("id"),
		})
	})

	_ = r.Start(context.Background())
}
