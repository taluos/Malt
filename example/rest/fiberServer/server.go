package fiberserver

import (
	"context"

	"github.com/gofiber/fiber/v3"

	httpserver "github.com/taluos/Malt/server/rest"
	fiberServer "github.com/taluos/Malt/server/rest/rest-fiber"
)

func RestInit() httpserver.Server {
	srv := httpserver.NewServer("fiber",
		fiberServer.WithAddress("127.0.0.1:9090"),
	)

	srv.Handle("GET", "/hello", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
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
