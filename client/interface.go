package client

import (
	"context"
	"fmt"

	"github.com/taluos/Malt/client/rest"
	"github.com/taluos/Malt/client/rpc"
	"github.com/taluos/Malt/pkg/errors"
)

// Client 统一的客户端接口
type Client interface {
	// Type 返回客户端类型（rest/rpc）
	Type() string
	// Close 关闭客户端
	Close(ctx context.Context) error
}

// NewClient 统一的客户端创建函数
func NewClient(clientType string, config any) (Client, error) {
	switch clientType {
	case RESTClientType:
		if cfg, ok := config.(RESTConfig); ok {
			return newRESTClient(cfg.Type, cfg.BaseURL, cfg.Options...)
		}
		return nil, errors.New("[Client] invalid REST client config")
	case RPCClientType:
		if cfg, ok := config.(RPCConfig); ok {
			return newRPCClient(cfg.Method, cfg.Options...)
		}
		return nil, errors.New("[Client] invalid RPC client config")
	default:
		return nil, errors.New(fmt.Sprintf("[Client] unsupported client type: %s", clientType))
	}
}

// RESTConfig REST客户端配置
type RESTConfig struct {
	Type    string // "http" or "fasthttp"
	BaseURL string
	Options []rest.ClientOption
}

// RPCConfig RPC客户端配置
type RPCConfig struct {
	Method  string // "grpc"
	Options []rpc.ClientOptions
}
