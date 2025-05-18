package restserver

import (
	"context"

	httpserver "Malt/server/rest/httpServer"

	"github.com/gin-gonic/gin"
)

func RestInit() *httpserver.Server {
	srv := httpserver.NewServer(
		httpserver.WithPort(8080),
		// httpserver.WithMode("debug"),
		httpserver.WithTrans("zh"),
		httpserver.WithHealthz(true),
		httpserver.WithEnableProfiling(true),
		httpserver.WithMiddleware(gin.Recovery()),
	)

	srv.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return srv
}

func RestRun(srv *httpserver.Server, ctx context.Context) error {
	return srv.Start(ctx)
}

func RestStop(srv *httpserver.Server, ctx context.Context) error {
	return srv.Stop(ctx)
}
