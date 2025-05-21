package rpcclient

import (
	"time"

	metric "github.com/taluos/Malt/core/metrics"
	"github.com/taluos/Malt/core/registry"
	maltAgent "github.com/taluos/Malt/core/trace"

	"google.golang.org/grpc"
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

// ClientOptions defines the method to customize a clientOptions.
type ClientOptions func(c *clientOptions)

func WithEndpoint(endpoint string) ClientOptions {
	return func(c *clientOptions) {
		c.endpoint = endpoint
	}
}

func WithTimeout(timeout time.Duration) ClientOptions {
	return func(c *clientOptions) {
		c.timeout = timeout
	}
}

func WithInsecure(insecure bool) ClientOptions {
	return func(c *clientOptions) {
		c.insecure = insecure
	}
}

func WithEnableTracing(enableTracing bool) ClientOptions {
	return func(c *clientOptions) {
		c.enableTracing = enableTracing
	}
}

func WithEnableMetrics(enableMetrics bool) ClientOptions {
	return func(c *clientOptions) {
		c.enableMetrics = enableMetrics
	}
}

func WithHistogramVecOpts(opts *metric.HistogramVecOpts) ClientOptions {
	return func(c *clientOptions) {
		c.histogramVecOpts = opts
	}
}

func WithCounterVecOpts(opts *metric.CounterVecOpts) ClientOptions {
	return func(c *clientOptions) {
		c.counterVecOpts = opts
	}
}

func WithDiscovery(discovery registry.Discovery) ClientOptions {
	return func(c *clientOptions) {
		c.discovery = discovery
	}
}

func WithAgent(agent *maltAgent.Agent) ClientOptions {
	return func(c *clientOptions) {
		c.agent = agent
	}
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) ClientOptions {
	return func(c *clientOptions) {
		c.unaryInterceptors = interceptors
	}
}

func WithStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) ClientOptions {
	return func(c *clientOptions) {
		c.streamInterceptors = interceptors
	}
}

func WithOptions(opts ...grpc.DialOption) ClientOptions {
	return func(c *clientOptions) {
		c.grpcOpts = opts
	}
}

func WithBalancerName(name string) ClientOptions {
	return func(c *clientOptions) {
		c.balancerName = name
	}
}
