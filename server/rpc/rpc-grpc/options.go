package grpc

import (
	"net"
	"net/url"
	"time"

	"github.com/taluos/Malt/api/metadata"
	metric "github.com/taluos/Malt/core/metrics"
	maltAgent "github.com/taluos/Malt/core/trace"
	auth "github.com/taluos/Malt/pkg/auth-jwt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

type serverOptions struct {
	name     string        // 服务器名称
	address  string        `validate:"required"`       // 服务器监听地址: ip:port
	endpoint *url.URL      `validate:"required"`       // 服务器端点URL: grpc://ip:port
	timeout  time.Duration `validate:"required,gte=0"` // 超时时间

	enableTracing     bool `validate:"required"` // 是否启用追踪
	enableMetrics     bool `validate:"required"` // 是否启用指标
	enableHealthCheck bool `validate:"required"` // 是否启用健康检查
	enableReflection  bool `validate:"required"` // 是否启用反射
	enableInsecure    bool `validate:"required"` // 是否启用不安全连接

	histogramVecOpts *metric.HistogramVecOpts
	counterVecOpts   *metric.CounterVecOpts

	unaryInterceptors  []grpc.UnaryServerInterceptor  // 一元拦截器列表
	streamInterceptors []grpc.StreamServerInterceptor // 流式拦截器列表
	grpcOpts           []grpc.ServerOption            // gRPC服务器选项

	listener    net.Listener     // 网络监听器
	metadata    *metadata.Server // 元数据服务器
	healthCheck *health.Server   // 健康检查服务器

	JWTauthenticator *auth.Authenticator // 认证器
	agent            *maltAgent.Agent
}

func (o *serverOptions) Validate() error {
	validator := validator.New()
	err := validator.Struct(o)
	if err != nil {
		return err
	}
	return nil
}

// ServerOptions 定义了自定义服务器选项的方法
type ServerOptions func(s *serverOptions)

func WithName(name string) ServerOptions {
	return func(s *serverOptions) {
		s.name = name
	}
}

func WithAuthenticator(authenticator *auth.Authenticator) ServerOptions {
	return func(s *serverOptions) {
		s.JWTauthenticator = authenticator
	}
}

func WithAgent(agent *maltAgent.Agent) ServerOptions {
	return func(s *serverOptions) {
		s.agent = agent
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

func WithEnableReflection(enableReflection bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableReflection = enableReflection
	}
}

func WithMetadata(metadata *metadata.Server) ServerOptions {
	return func(s *serverOptions) {
		s.metadata = metadata
	}
}
func WithEnableInsecure(enableInsecure bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableInsecure = enableInsecure
	}
}
