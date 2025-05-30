package restserver

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/taluos/Malt/server"
	httpserver "github.com/taluos/Malt/server/rest"
	ginServer "github.com/taluos/Malt/server/rest/rest-gin"
)

func RestInit() httpserver.Server {
	ServerInstance1, _ := server.NewServer("rest",
		server.RESTConfig{
			Method: "gin",
			Options: []httpserver.ServerOptions{
				ginServer.WithName("Malt Rest Server"),
				ginServer.WithHealthz(true),
				ginServer.WithAddress("127.0.0.1:8080"),
				ginServer.WithMiddleware(gin.Recovery()),
			},
		},
	)

	srv := ServerInstance1.(httpserver.Server)
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
