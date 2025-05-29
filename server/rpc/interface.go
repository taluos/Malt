package rpc

import (
	"context"
	"net/url"

	"github.com/taluos/Malt/pkg/log"
)

// Server 定义了RPC服务器的基本接口
type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error

	// Stop 停止服务器
	Stop(ctx context.Context) error

	// Endpoint 返回服务器的端点URL
	Endpoint() (*url.URL, error)

	// Engine 返回底层的RPC引擎，允许直接访问底层实现
	// 注意：这可能会导致与底层实现的耦合，应谨慎使用
	Engine() any

	// RegisterService 注册一个服务到RPC服务器
	RegisterService(desc any, impl any) Server
}

// ServerOptions 定义了创建服务器的选项
type ServerOptions interface{}

// NewServer 创建一个新的RPC服务器实例
func NewServer(method string, opts ...ServerOptions) Server {
	// 这里可以根据配置选择不同的实现
	switch method {
	case grpcServerType:
		return newGrpcServer(opts...)
	default:
		log.Errorf("[RPC-Server] unknown rpc server type: %s", method)
		return nil
	}
}
