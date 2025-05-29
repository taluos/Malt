package server

import (
	"context"
	"fmt"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/log"
	"github.com/taluos/Malt/server/rest"
	"github.com/taluos/Malt/server/rpc"
)

// Server 统一的服务器接口
type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error
	// Stop 停止服务器
	Stop(ctx context.Context) error
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
