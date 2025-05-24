package rest

import (
	"context"

	"github.com/gofiber/fiber/v3"
	httpserver "github.com/taluos/Malt/server/rest/rest-fiber"
)

type fiberServer struct {
	app *httpserver.Server
}

type fiberRouteGroup struct {
	group fiber.Router
}

var _ Server = (*fiberServer)(nil)
var _ RouteGroup = (*fiberRouteGroup)(nil)

func newFiberServer(opts ...ServerOptions) Server {
	serverOpts := convertFiberOptions(opts...)

	// 创建服务器
	app := httpserver.NewServer(serverOpts...)

	server := &fiberServer{
		app,
	}
	return server
}

func (s *fiberServer) Start(ctx context.Context) error {
	return s.app.Start(ctx)
}

func (s *fiberServer) Stop(ctx context.Context) error {
	return s.app.Stop(ctx)
}

func (s *fiberServer) App() any {
	return s.app.App
}

func (s *fiberServer) Group(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	group := s.app.Group(relativePath, fiberHandlers...)
	return &fiberRouteGroup{
		group: group,
	}
}

func (s *fiberServer) Use(middlewares ...any) Server {
	fiberMiddlewares := convertToFiberHandlers(middlewares...)
	for _, mw := range fiberMiddlewares {
		s.app.Use(mw)
	}
	return s
}

func (s *fiberServer) Handle(httpMethod, relativePath string, handlers ...any) Server {
	fiberHandlers := convertToFiberHandlers(handlers...)
	s.app.Add([]string{httpMethod}, relativePath, nil, fiberHandlers...)
	return s
}

// Group 实现RouteGroup.Group
func (g *fiberRouteGroup) Group(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	group := g.group.Group(relativePath, fiberHandlers...)
	return &fiberRouteGroup{group: group}
}

// Use 实现RouteGroup.Use
func (g *fiberRouteGroup) Use(middleware ...any) RouteGroup {
	fiberMiddleware := convertToFiberHandlers(middleware...)
	for _, mw := range fiberMiddleware {
		g.group.Use(mw)
	}
	return g
}

// Handle 实现RouteGroup.Handle
func (g *fiberRouteGroup) Handle(httpMethod, relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Add([]string{httpMethod}, relativePath, nil, fiberHandlers...)
	return g
}

// GET 实现RouteGroup.GET
func (g *fiberRouteGroup) GET(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Get(relativePath, nil, fiberHandlers...)
	return g
}

// POST 实现RouteGroup.POST
func (g *fiberRouteGroup) POST(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Post(relativePath, nil, fiberHandlers...)
	return g
}

// PUT 实现RouteGroup.PUT
func (g *fiberRouteGroup) PUT(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Put(relativePath, nil, fiberHandlers...)
	return g
}

// DELETE 实现RouteGroup.DELETE
func (g *fiberRouteGroup) DELETE(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Delete(relativePath, nil, fiberHandlers...)
	return g
}

// PATCH 实现RouteGroup.PATCH
func (g *fiberRouteGroup) PATCH(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Patch(relativePath, nil, fiberHandlers...)
	return g
}

// HEAD 实现RouteGroup.HEAD
func (g *fiberRouteGroup) HEAD(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Head(relativePath, nil, fiberHandlers...)
	return g
}

// OPTIONS 实现RouteGroup.OPTIONS
func (g *fiberRouteGroup) OPTIONS(relativePath string, handlers ...any) RouteGroup {
	fiberHandlers := convertToFiberHandlers(handlers...)
	g.group.Options(relativePath, nil, fiberHandlers...)
	return g
}

// 辅助函数：转换通用处理器为Fiber处理器
func convertToFiberHandlers(handlers ...any) []fiber.Handler {
	fiberHandlers := make([]fiber.Handler, 0, len(handlers))
	for _, h := range handlers {
		if fh, ok := h.(fiber.Handler); ok {
			fiberHandlers = append(fiberHandlers, fh)
		} else if fn, ok := h.(func(fiber.Ctx) error); ok {
			// 将函数转换为fiber.Handler
			fiberHandlers = append(fiberHandlers, fiber.Handler(fn))
		} else if fn, ok := h.(func(c fiber.Ctx) error); ok {
			// 另一种常见的函数签名
			fiberHandlers = append(fiberHandlers, fiber.Handler(fn))
		}
		// 可以添加更多类型的转换
	}
	return fiberHandlers
}

// 辅助函数：转换通用选项为Fiber选项
func convertFiberOptions(opts ...ServerOptions) []httpserver.ServerOptions {
	serverOpts := make([]httpserver.ServerOptions, 0, len(opts))
	for _, opt := range opts {
		if so, ok := opt.(httpserver.ServerOptions); ok {
			serverOpts = append(serverOpts, so)
		}
	}
	return serverOpts
}
