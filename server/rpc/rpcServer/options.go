package rpcserver

import (
	"net"
	"net/url"
	"time"

	"github.com/taluos/Malt/api/metadata"
	metric "github.com/taluos/Malt/core/metrics"
	"github.com/taluos/Malt/server/rpc/internal/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

type serverOptions struct {
	address  string        // 服务器监听地址: ip:port
	endpoint *url.URL      // 服务器端点URL: grpc://ip:port
	timeout  time.Duration // 超时时间

	enableTracing     bool // 是否启用追踪
	enableMetrics     bool // 是否启用指标
	enableHealthCheck bool // 是否启用健康检查
	enableReflection  bool // 是否启用反射

	histogramVecOpts *metric.HistogramVecOpts
	counterVecOpts   *metric.CounterVecOpts

	unaryInterceptors  []grpc.UnaryServerInterceptor  // 一元拦截器列表
	streamInterceptors []grpc.StreamServerInterceptor // 流式拦截器列表
	grpcOpts           []grpc.ServerOption            // gRPC服务器选项

	listener    net.Listener     // 网络监听器
	metadata    *metadata.Server // 元数据服务器
	healthCheck *health.Server   // 健康检查服务器

	authenticator *auth.Authenticator // 认证器

}

type ServerOptions func(s *serverOptions)

func WithAuthenticator(authenticator *auth.Authenticator) ServerOptions {
	return func(s *serverOptions) {
		s.authenticator = authenticator
	}
}

func WithAddress(address string) ServerOptions {
	return func(s *serverOptions) {
		s.address = address
	}
}

func WithEndpoint(endpoint *url.URL) ServerOptions {
	return func(s *serverOptions) {
		s.endpoint = endpoint
	}
}

func WithListener(listener net.Listener) ServerOptions {
	return func(s *serverOptions) {
		s.listener = listener
	}
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOptions {
	return func(s *serverOptions) {
		s.unaryInterceptors = interceptors
	}
}

func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOptions {
	return func(s *serverOptions) {
		s.streamInterceptors = interceptors
	}
}

func WithOptions(opts ...grpc.ServerOption) ServerOptions {
	return func(s *serverOptions) {
		s.grpcOpts = opts
	}
}

func WithHealthCheck(healthCheck *health.Server) ServerOptions {
	return func(s *serverOptions) {
		s.healthCheck = healthCheck
	}
}

func WithTimeout(timeout time.Duration) ServerOptions {
	return func(s *serverOptions) {
		s.timeout = timeout
	}
}

func WithEnableTracing(enableTracing bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableTracing = enableTracing
	}
}

func WithEnableMetrics(enableMetrics bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableMetrics = enableMetrics
	}
}

func WithHistogramVecOpts(opts *metric.HistogramVecOpts) ServerOptions {
	return func(s *serverOptions) {
		s.histogramVecOpts = opts
	}
}

func WithCounterVecOpts(opts *metric.CounterVecOpts) ServerOptions {
	return func(s *serverOptions) {
		s.counterVecOpts = opts
	}
}

func WithEnableHealthCheck(enableHealthCheck bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableHealthCheck = enableHealthCheck
	}
}

func WithMetadata(metadata *metadata.Server) ServerOptions {
	return func(s *serverOptions) {
		s.metadata = metadata
	}
}
