package rpc

import (
	"context"
	"log"
	"net/url"
	"reflect"

	gRpc "github.com/taluos/Malt/server/rpc/rpc-grpc"

	"google.golang.org/grpc"
)

// grpcServer 是基于gRPC的Server实现
type grpcServer struct {
	server *gRpc.Server
}

// grpcClient 是基于gRPC的Client实现
type grpcClient struct {
	client *gRpc.Client
}

// 确保grpcServer实现了Server接口
var _ Server = (*grpcServer)(nil)

// 确保grpcClient实现了Client接口
var _ Client = (*grpcClient)(nil)

// newGrpcServer 创建一个新的基于gRPC的服务器
func newGrpcServer(opts ...ServerOptions) Server {
	// 转换选项
	serverOpts := convertOptions(opts...)

	// 创建服务器
	server := gRpc.NewServer(serverOpts...)

	return &grpcServer{
		server: server,
	}
}

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

// Start 实现Server.Start
func (s *grpcServer) Start(ctx context.Context) error {
	return s.server.Start(ctx)
}

// Stop 实现Server.Stop
func (s *grpcServer) Stop(ctx context.Context) error {
	return s.server.Stop(ctx)
}

// Endpoint 实现Server.Endpoint
func (s *grpcServer) Endpoint() (*url.URL, error) {
	return s.server.Endpoint()
}

// Engine 实现Server.Engine
func (s *grpcServer) Engine() any {
	return s.server.Server
}

// RegisterService 实现Server.RegisterService
func (s *grpcServer) RegisterService(desc interface{}, impl interface{}) Server {
	// 这里需要根据gRPC的注册方式进行适配
	// 例如：desc 可能是 *grpc.ServiceDesc，impl 是服务实现
	if sd, ok := desc.(*grpc.ServiceDesc); ok {
		s.server.RegisterService(sd, impl)
	} else if registerFunc, ok := desc.(func(s grpc.ServiceRegistrar, srv interface{})); ok {
		// 支持 protobuf 生成的注册函数
		registerFunc(s.server.Server, impl)
	} else {
		// 尝试通过反射调用注册函数
		registerFuncValue := reflect.ValueOf(desc)
		if registerFuncValue.Kind() == reflect.Func &&
			registerFuncValue.Type().NumIn() == 2 &&
			registerFuncValue.Type().In(0).Implements(reflect.TypeOf((*grpc.ServiceRegistrar)(nil)).Elem()) {

			args := []reflect.Value{
				reflect.ValueOf(s.server.Server),
				reflect.ValueOf(impl),
			}
			registerFuncValue.Call(args)
		} else {
			// 记录错误日志
			log.Printf("无法注册服务，未知的注册函数类型: %T", desc)
		}
	}
	return s
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

// 辅助函数：转换通用选项为gRPC选项
func convertOptions(opts ...ServerOptions) []gRpc.ServerOptions {
	serverOpts := make([]gRpc.ServerOptions, 0, len(opts))
	for _, opt := range opts {
		if so, ok := opt.(gRpc.ServerOptions); ok {
			serverOpts = append(serverOpts, so)
		}
	}
	return serverOpts
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
