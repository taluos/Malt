package grpc

import (
	"net"
	"net/url"
	"time"

	"github.com/taluos/Malt/api/metadata"
	metric "github.com/taluos/Malt/core/metrics"
	"github.com/taluos/Malt/core/registry"
	maltAgent "github.com/taluos/Malt/core/trace"
	"github.com/taluos/Malt/pkg/auth-jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

// A clientOptions is a client option.
type clientOptions struct {
	endpoint string // 服务器地址
	timeout  time.Duration

	insecure      bool
	enableTracing bool
	enableMetrics bool

	histogramVecOpts *metric.HistogramVecOpts
	counterVecOpts   *metric.CounterVecOpts

	discovery registry.Discovery
	agent     *maltAgent.Agent

	unaryInterceptors  []grpc.UnaryClientInterceptor  // 一元拦截器列表
	streamInterceptors []grpc.StreamClientInterceptor // 流式拦截器列表
	grpcOpts           []grpc.DialOption

	balancerName string // 负载均衡器名称
}

type serverOptions struct {
	address  string        // 服务器监听地址: ip:port
	endpoint *url.URL      // 服务器端点URL: grpc://ip:port
	timeout  time.Duration // 超时时间

	enableTracing     bool // 是否启用追踪
	enableMetrics     bool // 是否启用指标
	enableHealthCheck bool // 是否启用健康检查
	enableReflection  bool // 是否启用反射
	enableInsecure    bool // 是否启用不安全连接

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

// ClientOptions 定义了自定义客户端选项的方法
type ClientOptions func(c *clientOptions)

// ServerOptions 定义了自定义服务器选项的方法
type ServerOptions func(s *serverOptions)

func WithClientEndpoint(endpoint string) ClientOptions {
	return func(c *clientOptions) {
		c.endpoint = endpoint
	}
}

func WithClientTimeout(timeout time.Duration) ClientOptions {
	return func(c *clientOptions) {
		c.timeout = timeout
	}
}

func WithClientInsecure(insecure bool) ClientOptions {
	return func(c *clientOptions) {
		c.insecure = insecure
	}
}

func WithClientEnableTracing(enableTracing bool) ClientOptions {
	return func(c *clientOptions) {
		c.enableTracing = enableTracing
	}
}

func WithClientEnableMetrics(enableMetrics bool) ClientOptions {
	return func(c *clientOptions) {
		c.enableMetrics = enableMetrics
	}
}

func WithClientHistogramVecOpts(opts *metric.HistogramVecOpts) ClientOptions {
	return func(c *clientOptions) {
		c.histogramVecOpts = opts
	}
}

func WithClientCounterVecOpts(opts *metric.CounterVecOpts) ClientOptions {
	return func(c *clientOptions) {
		c.counterVecOpts = opts
	}
}

func WithClientDiscovery(discovery registry.Discovery) ClientOptions {
	return func(c *clientOptions) {
		c.discovery = discovery
	}
}

func WithClientAgent(agent *maltAgent.Agent) ClientOptions {
	return func(c *clientOptions) {
		c.agent = agent
	}
}

func WithClientUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) ClientOptions {
	return func(c *clientOptions) {
		c.unaryInterceptors = interceptors
	}
}

func WithClientStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) ClientOptions {
	return func(c *clientOptions) {
		c.streamInterceptors = interceptors
	}
}

func WithClientOptions(opts ...grpc.DialOption) ClientOptions {
	return func(c *clientOptions) {
		c.grpcOpts = opts
	}
}

func WithClientBalancerName(name string) ClientOptions {
	return func(c *clientOptions) {
		c.balancerName = name
	}
}

func WithServerAuthenticator(authenticator *auth.Authenticator) ServerOptions {
	return func(s *serverOptions) {
		s.JWTauthenticator = authenticator
	}
}

func WithServerAgent(agent *maltAgent.Agent) ServerOptions {
	return func(s *serverOptions) {
		s.agent = agent
	}
}

func WithServerAddress(address string) ServerOptions {
	return func(s *serverOptions) {
		s.address = address
	}
}

func WithServerEndpoint(endpoint *url.URL) ServerOptions {
	return func(s *serverOptions) {
		s.endpoint = endpoint
	}
}

func WithServerListener(listener net.Listener) ServerOptions {
	return func(s *serverOptions) {
		s.listener = listener
	}
}

func WithServerUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOptions {
	return func(s *serverOptions) {
		s.unaryInterceptors = interceptors
	}
}

func WithServerStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOptions {
	return func(s *serverOptions) {
		s.streamInterceptors = interceptors
	}
}

func WithServerOptions(opts ...grpc.ServerOption) ServerOptions {
	return func(s *serverOptions) {
		s.grpcOpts = opts
	}
}

func WithServerHealthCheck(healthCheck *health.Server) ServerOptions {
	return func(s *serverOptions) {
		s.healthCheck = healthCheck
	}
}

func WithServerTimeout(timeout time.Duration) ServerOptions {
	return func(s *serverOptions) {
		s.timeout = timeout
	}
}

func WithServerEnableTracing(enableTracing bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableTracing = enableTracing
	}
}

func WithServerEnableMetrics(enableMetrics bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableMetrics = enableMetrics
	}
}

func WithServerHistogramVecOpts(opts *metric.HistogramVecOpts) ServerOptions {
	return func(s *serverOptions) {
		s.histogramVecOpts = opts
	}
}

func WithServerCounterVecOpts(opts *metric.CounterVecOpts) ServerOptions {
	return func(s *serverOptions) {
		s.counterVecOpts = opts
	}
}

func WithServerEnableHealthCheck(enableHealthCheck bool) ServerOptions {
	return func(s *serverOptions) {
		s.enableHealthCheck = enableHealthCheck
	}
}

func WithServerMetadata(metadata *metadata.Server) ServerOptions {
	return func(s *serverOptions) {
		s.metadata = metadata
	}
}
