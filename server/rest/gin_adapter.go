package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	httpserver "github.com/taluos/Malt/server/rest/rest-gin"
)

// ginServer 是基于Gin的Server实现
type ginServer struct {
	server *httpserver.Server
}

// ginRouteGroup 是基于Gin的RouteGroup实现
type ginRouteGroup struct {
	group *gin.RouterGroup
}

// 确保ginServer实现了Server接口
var _ Server = (*ginServer)(nil)

// 确保ginRouteGroup实现了RouteGroup接口
var _ RouteGroup = (*ginRouteGroup)(nil)

// newGinServer 创建一个新的基于Gin的服务器
func newGinServer(opts ...ServerOptions) Server {
	// 转换选项
	serverOpts := convertGinOptions(opts...)

	// 创建服务器
	server := httpserver.NewServer(serverOpts...)

	return &ginServer{
		server: server,
	}
}

// Start 实现Server.Start
func (s *ginServer) Start(ctx context.Context) error {
	return s.server.Start(ctx)
}

// Stop 实现Server.Stop
func (s *ginServer) Stop(ctx context.Context) error {
	return s.server.Stop(ctx)
}

// Engine 实现Server.Engine
func (s *ginServer) Engine() any {
	return s.server.Engine
}

// Group 实现Server.Group
func (s *ginServer) Group(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	group := s.server.Group(relativePath, ginHandlers...)
	return &ginRouteGroup{group: group}
}

// Use 实现Server.Use
func (s *ginServer) Use(middleware ...any) Server {
	ginMiddleware := convertToGinHandlers(middleware...)
	s.server.Use(ginMiddleware...)
	return s
}

// Handle 实现Server.Handle
func (s *ginServer) Handle(httpMethod, relativePath string, handlers ...any) Server {
	ginHandlers := convertToGinHandlers(handlers...)
	s.server.Handle(httpMethod, relativePath, ginHandlers...)
	return s
}

// Group 实现RouteGroup.Group
func (g *ginRouteGroup) Group(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	group := g.group.Group(relativePath, ginHandlers...)
	return &ginRouteGroup{group: group}
}

// Use 实现RouteGroup.Use
func (g *ginRouteGroup) Use(middleware ...any) RouteGroup {
	ginMiddleware := convertToGinHandlers(middleware...)
	g.group.Use(ginMiddleware...)
	return g
}

// Handle 实现RouteGroup.Handle
func (g *ginRouteGroup) Handle(httpMethod, relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.Handle(httpMethod, relativePath, ginHandlers...)
	return g
}

// GET 实现RouteGroup.GET
func (g *ginRouteGroup) GET(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.GET(relativePath, ginHandlers...)
	return g
}

// POST 实现RouteGroup.POST
func (g *ginRouteGroup) POST(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.POST(relativePath, ginHandlers...)
	return g
}

// PUT 实现RouteGroup.PUT
func (g *ginRouteGroup) PUT(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.PUT(relativePath, ginHandlers...)
	return g
}

// DELETE 实现RouteGroup.DELETE
func (g *ginRouteGroup) DELETE(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.DELETE(relativePath, ginHandlers...)
	return g
}

// PATCH 实现RouteGroup.PATCH
func (g *ginRouteGroup) PATCH(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.PATCH(relativePath, ginHandlers...)
	return g
}

// HEAD 实现RouteGroup.HEAD
func (g *ginRouteGroup) HEAD(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.HEAD(relativePath, ginHandlers...)
	return g
}

// OPTIONS 实现RouteGroup.OPTIONS
func (g *ginRouteGroup) OPTIONS(relativePath string, handlers ...any) RouteGroup {
	ginHandlers := convertToGinHandlers(handlers...)
	g.group.OPTIONS(relativePath, ginHandlers...)
	return g
}

// 辅助函数：转换通用处理器为Gin处理器
func convertToGinHandlers(handlers ...any) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		if gh, ok := h.(gin.HandlerFunc); ok {
			ginHandlers = append(ginHandlers, gh)
		} else if fn, ok := h.(func(*gin.Context)); ok {
			// 将函数转换为gin.HandlerFunc
			ginHandlers = append(ginHandlers, gin.HandlerFunc(fn))
		} else if fn, ok := h.(func(c *gin.Context)); ok {
			// 另一种常见的函数签名
			ginHandlers = append(ginHandlers, gin.HandlerFunc(fn))
		}
		// 可以添加更多类型的转换
	}
	return ginHandlers
}

// 辅助函数：转换通用选项为Gin选项
func convertGinOptions(opts ...ServerOptions) []httpserver.ServerOptions {
	serverOpts := make([]httpserver.ServerOptions, 0, len(opts))
	for _, opt := range opts {
		if so, ok := opt.(httpserver.ServerOptions); ok {
			serverOpts = append(serverOpts, so)
		}
	}
	return serverOpts
}
