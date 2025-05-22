package restserver

import (
	"context"

	"github.com/gin-gonic/gin"
	httpserver "github.com/taluos/Malt/server/rest"
	ginServer "github.com/taluos/Malt/server/rest/rest-gin"
)

func RestInit() httpserver.Server {

	srv := httpserver.NewServer("gin",
		ginServer.WithPort(8080),
		ginServer.WithMode("debug"),
	)

	srv.Handle("GET", "/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	InitRouter(srv)

	return srv
}

func RestRun(srv httpserver.Server, ctx context.Context) error {
	return srv.Start(ctx)
}

func RestStop(srv httpserver.Server, ctx context.Context) error {
	return srv.Stop(ctx)
}
