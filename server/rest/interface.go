package rest

import (
	"context"
)

// Server 定义了REST服务器的基本接口
type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error

	// Stop 停止服务器
	Stop(ctx context.Context) error

	// Engine 返回底层的HTTP引擎，允许直接访问底层实现
	// 注意：这可能会导致与底层实现的耦合，应谨慎使用
	Engine() any

	// Group 创建一个新的路由组
	Group(relativePath string, handlers ...any) RouteGroup

	// Use 添加中间件
	Use(middleware ...any) Server

	// Handle 注册一个新路由
	Handle(httpMethod, relativePath string, handlers ...any) Server
}

// RouteGroup 定义了路由组的接口
type RouteGroup interface {
	// Group 创建一个嵌套的路由组
	Group(relativePath string, handlers ...interface{}) RouteGroup

	// Use 为路由组添加中间件
	Use(middleware ...interface{}) RouteGroup

	// Handle 在路由组中注册一个新路由
	Handle(httpMethod, relativePath string, handlers ...interface{}) RouteGroup

	// GET 注册一个GET请求处理器
	GET(relativePath string, handlers ...interface{}) RouteGroup

	// POST 注册一个POST请求处理器
	POST(relativePath string, handlers ...interface{}) RouteGroup

	// PUT 注册一个PUT请求处理器
	PUT(relativePath string, handlers ...interface{}) RouteGroup

	// DELETE 注册一个DELETE请求处理器
	DELETE(relativePath string, handlers ...interface{}) RouteGroup

	// PATCH 注册一个PATCH请求处理器
	PATCH(relativePath string, handlers ...interface{}) RouteGroup

	// HEAD 注册一个HEAD请求处理器
	HEAD(relativePath string, handlers ...interface{}) RouteGroup

	// OPTIONS 注册一个OPTIONS请求处理器
	OPTIONS(relativePath string, handlers ...interface{}) RouteGroup
}

// ServerOptions 定义了创建服务器的选项
type ServerOptions any

// NewServer 创建一个新的REST服务器实例
func NewServer(method string, opts ...ServerOptions) Server {
	factory, exists := serverFactories[method]
	if !exists {
		return nil
	}
	return factory(opts...)
}
