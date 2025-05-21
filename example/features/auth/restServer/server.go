package restserver

import (
	"context"

	httpserver "github.com/taluos/Malt/server/rest/Server"

	"github.com/gin-gonic/gin"
)

func RestInit() *httpserver.Server {
	srv := httpserver.NewServer(
		httpserver.WithPort(8080),
		httpserver.WithMode("debug"),
		httpserver.WithHealthz(true),
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
