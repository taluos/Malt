package rpc

import (
	"context"

	gRpc "github.com/taluos/Malt/client/rpc/rpc-grpc"
)

// grpcClient 是基于gRPC的Client实现
type grpcClient struct {
	client *gRpc.Client
}

// 确保grpcClient实现了Client接口
var _ Client = (*grpcClient)(nil)

// newGrpcClient 创建一个新的基于gRPC的客户端
func newGrpcClient(opts ...ClientOptions) (Client, error) {
	// 转换选项
	clientOpts := convertClientOptions(opts...)
	// 创建客户端
	client, err := gRpc.NewClient(clientOpts...)
	if err != nil {
		return nil, err
	}
	return &grpcClient{
		client: client,
	}, nil
}

// Endpoint 实现Client.Endpoint
func (c *grpcClient) Endpoint() string {
	return c.client.Endpoint()
}

// Close 实现Client.Close
func (c *grpcClient) Close(ctx context.Context) error {
	return c.client.Close(ctx)
}

// Conn 实现Client.Conn
func (c *grpcClient) Conn() any {
	return c.client.ClientConn
}

// 辅助函数：转换通用选项为gRPC客户端选项
func convertClientOptions(opts ...ClientOptions) []gRpc.ClientOptions {
	clientOpts := make([]gRpc.ClientOptions, 0, len(opts))
	for _, opt := range opts {
		if co, ok := opt.(gRpc.ClientOptions); ok {
			clientOpts = append(clientOpts, co)
		}
	}
	return clientOpts
}
