package client

import (
	"context"

	"github.com/taluos/Malt/client/rpc"
)

// rpcClientWrapper RPC客户端包装器
type rpcClientWrapper struct {
	client rpc.Client
}

// RPCClient RPC客户端接口，继承统一接口
type RPCClient interface {
	Type() string
	rpc.Client // 嵌入RPC客户端接口
}

var _ Client = (*restClientWrapper)(nil)

// NewRPCClient 创建RPC客户端
func newRPCClient(method string, opts ...rpc.ClientOptions) (RPCClient, error) {
	client, err := rpc.NewClient(method, opts...)
	if err != nil {
		return nil, err
	}
	return &rpcClientWrapper{client: client}, nil
}

func (w *rpcClientWrapper) Type() string {
	return RPCClientType
}

func (w *rpcClientWrapper) Close(ctx context.Context) error {
	return w.client.Close(ctx)
}

func (w *rpcClientWrapper) Endpoint() string {
	return w.client.Endpoint()
}

func (w *rpcClientWrapper) Conn() any {
	return w.client.Conn()
}
