package server

import (
	"context"
	"fmt"
	"net/url"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/log"
	"github.com/taluos/Malt/server/rest"
	"github.com/taluos/Malt/server/rpc"
)

var _ Server = (RESTServer)(nil)
var _ Server = (RPCServer)(nil)

// Server 统一的服务器接口
type Server interface {
	// Type 服务器类型
	Type() string
	// Start 启动服务器
	Start(ctx context.Context) error
	// Stop 停止服务器
	Stop(ctx context.Context) error
}

// RESTServer 扩展REST服务器接口
type RESTServer interface {
	Server
	Group(relativePath string, handlers ...any) rest.RouteGroup
	Use(middleware ...any) rest.Server
	Handle(httpMethod, relativePath string, handlers ...any) rest.Server
}

// RPCServer 扩展RPC服务器接口
type RPCServer interface {
	Server
	Endpoint() (*url.URL, error)
	Engine() any
	RegisterService(desc any, impl any) rpc.Server
}

// NewServer 统一的服务器创建函数
func NewServer(serverType string, config any) (Server, error) {
	switch serverType {
	case RESTServerType:
		if cfg, ok := config.(RESTConfig); ok {
			return newRestServer(cfg.Method, cfg.Options...)
		}
		return nil, errors.New("[Server] invalid REST server config")
	case RPCServerType:
		if cfg, ok := config.(RPCConfig); ok {
			return newRPCServer(cfg.Method, cfg.Options...)
		}
		return nil, errors.New("[Server] invalid RPC server config")
	default:
		log.Errorf("[Server] unsupported server type: %s", serverType)
		return nil, errors.New(fmt.Sprintf("[Server] unsupported server type: %s", serverType))
	}
}

// RESTConfig REST服务器配置
type RESTConfig struct {
	Method  string // "gin", "fiber", etc.
	Options []rest.ServerOptions
}

// RPCConfig RPC服务器配置
type RPCConfig struct {
	Method  string // "grpc", etc.
	Options []rpc.ServerOptions
}
