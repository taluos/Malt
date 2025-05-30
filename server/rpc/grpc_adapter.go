package rpc

import (
	"context"
	"log"
	"net/url"
	"reflect"

	grpcServer "github.com/taluos/Malt/server/rpc/rpc-grpc"

	"google.golang.org/grpc"
)

// grpcServerWrapper 是基于gRPC的Server实现
type grpcServerWrapper struct {
	server *grpcServer.Server
}

// 确保grpcServer实现了Server接口
var _ Server = (*grpcServerWrapper)(nil)

// newGrpcServer 创建一个新的基于gRPC的服务器
func newGrpcServer(opts ...ServerOptions) *grpcServerWrapper {
	// 转换选项
	serverOpts := convertOptions(opts...)
	// 创建服务器
	server := grpcServer.NewServer(serverOpts...)
	return &grpcServerWrapper{
		server: server,
	}
}

func (s *grpcServerWrapper) Type() string {
	return "grpc"
}

// Start 实现Server.Start
func (s *grpcServerWrapper) Start(ctx context.Context) error {
	return s.server.Start(ctx)
}

// Stop 实现Server.Stop
func (s *grpcServerWrapper) Stop(ctx context.Context) error {
	return s.server.Stop(ctx)
}

// Endpoint 实现Server.Endpoint
func (s *grpcServerWrapper) Endpoint() (*url.URL, error) {
	return s.server.Endpoint()
}

// Engine 实现Server.Engine
func (s *grpcServerWrapper) Engine() any {
	return s.server.Server
}

// RegisterService 实现Server.RegisterService
func (s *grpcServerWrapper) RegisterService(desc interface{}, impl interface{}) Server {
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

// 辅助函数：转换通用选项为gRPC选项
func convertOptions(opts ...ServerOptions) []grpcServer.ServerOptions {
	serverOpts := make([]grpcServer.ServerOptions, 0, len(opts))
	for _, opt := range opts {
		if so, ok := opt.(grpcServer.ServerOptions); ok {
			serverOpts = append(serverOpts, so)
		}
	}
	return serverOpts
}
