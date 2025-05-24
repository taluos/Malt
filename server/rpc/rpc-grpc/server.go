package grpc

import (
	"context"
	"errors"
	"net"
	"net/url"

	"github.com/taluos/Malt/api/metadata"
	"github.com/taluos/Malt/pkg/host"
	"github.com/taluos/Malt/pkg/log"
	"github.com/taluos/Malt/server/rpc/rpc-grpc/internal/serverinterceptors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server 结构体定义了gRPC服务器的基本属性和配置
type Server struct {
	*grpc.Server // gRPC服务器实例
	rootCtx      context.Context
	opt          *serverOptions   // 服务器选项
	metadata     *metadata.Server // 元数据服务器
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

	if o.enableTracing && o.agent != nil {
		//grpcOptions = append(grpcOptions, grpc.StatsHandler(otelgrpc.NewServerHandler()))
		uraryInts = append(uraryInts,
			serverinterceptors.UnaryTracingInterceptor(o.agent))
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
	var err error

	if s.opt.address == "" {
		return errors.New("[gRPC] server address cannot be empty")
	}

	// 解析address：如果用户没有设置address，则从listener中获取一个address
	err = s.listenAndEndpoint()
	if err != nil {
		log.Errorf("[gRPC] Get endpoint failed: %s", err)
		return err
	}

	s.rootCtx = ctx

	log.Infof("[gRPC] server listening at %s", s.opt.address)

	s.opt.healthCheck.Resume()

	err = s.Server.Serve(s.opt.listener)
	if err != nil {
		log.Errorf("[gRPC] server serve failed: %s", err)
		return err
	}

	return err
}

func (s *Server) Stop(ctx context.Context) error {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

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
