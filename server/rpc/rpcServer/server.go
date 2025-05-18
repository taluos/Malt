package rpcserver

import (
	"Malt/api/metadata"
	"Malt/pkg/host"
	"Malt/pkg/log"
	"Malt/server/rpc/internal/serverinterceptors"
	"context"

	"net"
	"net/url"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server 结构体定义了gRPC服务器的基本属性和配置
type Server struct {
	*grpc.Server // gRPC服务器实例

	baseCtx context.Context

	opt *serverOptions // 服务器选项

	metadata *metadata.Server // 元数据服务器
}

type RPCServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Endpoint() (*url.URL, error)
	listenAndEndpoint() error
}

// NewServer 创建一个新的gRPC服务器实例
func NewServer(opts ...ServerOptions) *Server {

	o := &serverOptions{
		address:     defaultAddress,
		healthCheck: health.NewServer(),
		timeout:     defaultTimeout,

		enableTracing:     false,
		enableMetrics:     false,
		enableHealthCheck: true,
		enableReflection:  true,
	}

	for _, opt := range opts {
		opt(o)
	}

	uraryInts := []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryRecoverInterceptor,
		serverinterceptors.UnaryTimeoutInterceptor(o.timeout),
	}

	if o.enableMetrics {
		uraryInts = append(uraryInts,
			serverinterceptors.UnaryPrometheusInterceptor(o.histogramVecOpts, o.counterVecOpts))
	}

	if len(o.unaryInterceptors) > 0 {
		uraryInts = append(uraryInts, o.unaryInterceptors...)
	}

	streamInts := []grpc.StreamServerInterceptor{
		serverinterceptors.StreamRecoverInterceptor,
	}
	if len(o.streamInterceptors) > 0 {
		streamInts = append(streamInts, o.streamInterceptors...)
	}

	grpcOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(uraryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}

	// 将用户自己传入的grpc serverOptions合并到grpcOptions
	if len(o.grpcOpts) > 0 {
		grpcOptions = append(grpcOptions, o.grpcOpts...)
	}

	if o.enableTracing {
		grpcOptions = append(grpcOptions, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	// 删除或注释掉这段代码
	// if o.enableMetrics {
	//     grpcOptions = append(grpcOptions, grpc.UnaryInterceptor(
	//         serverinterceptors.UnaryPrometheusInterceptor(o.histogramVecOpts, o.counterVecOpts)))
	// }

	s := &Server{
		opt: o,
	}

	// 创建grpc server
	s.Server = grpc.NewServer(grpcOptions...)

	// 注册 metadata 服务
	s.metadata = metadata.NewServer(s.Server)
	metadata.RegisterMetadataServer(s.Server, s.metadata)

	// health register - 注册健康检查服务
	if s.opt.enableHealthCheck {
		grpc_health_v1.RegisterHealthServer(s.Server, s.opt.healthCheck)
	}

	//  注册反射服务，支持服务发现
	if s.opt.enableReflection {
		reflection.Register(s.Server)
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {

	// 解析address：如果用户没有设置address，则从listener中获取一个address
	err := s.listenAndEndpoint()
	if err != nil {
		log.Errorf("Get endpoint failed: %s", err)
		return err
	}

	s.baseCtx = ctx

	log.Infof("gRPC server listening at %s", s.opt.address)
	// 监听端口并启动服务
	s.opt.healthCheck.Resume()

	return s.Server.Serve(s.opt.listener)
}

func (s *Server) Stop(ctx context.Context) error {

	s.opt.healthCheck.Shutdown()

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Infof("[gRPC] server stopping")
		s.Server.GracefulStop()
	}()

	select {
	case <-done:
	case <-ctx.Done():
		log.Warn("[gRPC] server couldn't stop gracefully in time, doing force stop")
		s.Server.Stop()
	}

	log.Infof("[gRPC] server stopped")

	return nil
}

// Endpoint return a real address to registry endpoint.
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.opt.endpoint, nil
}

// ListenEndpoint 设置服务器的监听端点
// 如果没有传入address就从listener中获取一个endpoint
// ip 和 端口 的抽取
func (s *Server) listenAndEndpoint() error {

	// 如果用户已经设置了listener，则直接使用用户设置的listener
	if s.opt.listener == nil {
		lis, err := net.Listen("tcp", s.opt.address)
		if err != nil {
			log.Errorf("[gRPC] Listen to the listener failed: %s", err)
			return err
		}
		s.opt.listener = lis
	}

	// 提取地址
	address, err := host.Extract(s.opt.address, s.opt.listener)
	if err != nil {
		log.Errorf("[gRPC] Get address from listener failed: %s", err)
		closeErr := s.opt.listener.Close()
		if closeErr != nil {
			log.Errorf("[gRPC] Close listener failed: %s", closeErr)
			return closeErr
		}
		return err
	}

	s.opt.endpoint = &url.URL{
		Scheme: "grpc",
		Host:   address,
	}

	return nil
}
