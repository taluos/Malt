package server

import (
	"context"
	"fmt"
	"net/url"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/server/rpc"
)

// rpcServerWrapper RPC服务器包装器
type rpcServerWrapper struct {
	server     rpc.Server
	serverType string
}

var _ RPCServer = (*rpcServerWrapper)(nil)

// NewRPCServer 创建RPC服务器
func newRPCServer(serverType string, opts ...rpc.ServerOptions) (*rpcServerWrapper, error) {
	server := rpc.NewServer(serverType, opts...)
	if server == nil {
		return nil, errors.New(fmt.Sprintf("[Server] failed to create RPC server with method: %s", serverType))
	}
	return &rpcServerWrapper{server: server, serverType: serverType}, nil
}

func (w *rpcServerWrapper) Type() string {
	return w.serverType
}

func (w *rpcServerWrapper) Start(ctx context.Context) error {
	return w.server.Start(ctx)
}

func (w *rpcServerWrapper) Stop(ctx context.Context) error {
	return w.server.Stop(ctx)
}

func (w *rpcServerWrapper) Endpoint() (*url.URL, error) {
	return w.server.Endpoint()
}

func (w *rpcServerWrapper) Engine() any {
	return w.server.Engine()
}

func (w *rpcServerWrapper) RegisterService(desc any, impl any) rpc.Server {
	return w.server.RegisterService(desc, impl)
}
