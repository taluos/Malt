package rpc

import (
	"context"
	"fmt"

	"github.com/taluos/Malt/pkg/errors"
)

// Client 定义了RPC客户端的基本接口
type Client interface {
	// Endpoint 返回客户端连接的端点
	Endpoint() string
	// Close 关闭客户端连接
	Close(ctx context.Context) error
	// Conn 返回底层的RPC连接，允许直接访问底层实现
	// 注意：这可能会导致与底层实现的耦合，应谨慎使用
	Conn() any
}

// ClientOptions 定义了创建客户端的选项
type ClientOptions interface{}

const (
	GRPCClientType = "grpc"
)

// NewClient 创建一个新的RPC客户端实例
func NewClient(method string, opts ...ClientOptions) (Client, error) {
	// 这里可以根据配置选择不同的实现
	switch method {
	case GRPCClientType:
		return newGrpcClient(opts...)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported client type: %s", method))
	}
}
