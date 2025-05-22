package main

import (
	"context"
	"log"

	httpserver "github.com/taluos/Malt/server/rest"
	ginServer "github.com/taluos/Malt/server/rest/rest-gin"

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

	r := httpserver.NewServer("gin",
		ginServer.WithPort(8090),
		ginServer.WithEnableProfiling(true),
		ginServer.WithEnableTracing(true),
		ginServer.WithEnableMetrics(true),
		ginServer.WithMiddleware(gin.Recovery()),
	)

	r.Handle("GET", "/test/:id", func(c *gin.Context) {
		// 更新指标值
		_ = ginmetrics.GetMonitor().GetMetric("example_gauge_metric").Inc([]string{"label_value1"})

		c.JSON(200, map[string]string{
			"productId": c.Param("id"),
		})
	})

	_ = r.Start(context.Background())
}
