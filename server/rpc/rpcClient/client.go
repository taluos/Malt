package rpcclient

import (
	"context"

	"github.com/taluos/Malt/core/resolver/direct"
	"github.com/taluos/Malt/core/resolver/discovery"
	"github.com/taluos/Malt/pkg/log"
	"github.com/taluos/Malt/server/rpc/internal/clientinterceptors"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type (
	Client struct {
		// grpc.ClientConn
		CliConn *grpc.ClientConn

		opts clientOptions
	}

	// RPCClientMethod interface 定义了 RPC 客户端的基本方法
	RPCClientMethod interface {
		Endpoint() string
		Close(context.Context) error
	}
)

func NewClient(opts ...ClientOptions) (*Client, error) {

	o := clientOptions{
		endpoint: "/",

		insecure:      false,
		enableTracing: false,
		enableMetrics: false,

		timeout:      dialTimeout,
		balancerName: defautBalancer,
	}

	for _, opt := range opts {
		opt(&o)
	}

	CliConn, err := dial(o.insecure, o)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		CliConn: CliConn,
		opts:    o,
	}

	return cli, nil
}

func (c *Client) Endpoint() string {
	return c.opts.endpoint
}

func (c *Client) Close(ctx context.Context) error {

	if c.CliConn == nil {
		return nil
	}

	// 创建一个带超时的上下文
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, dialTimeout)
		defer cancel()
	}

	// 可以在这里添加其他资源的清理
	// 例如：等待所有正在进行的请求完成或超时

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Infof("[gRPC] client closing")
		c.CliConn.Close()
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return ctx.Err()
	}

	log.Infof("[gRPC] client closed")
	return nil
}

func dial(insecure bool, opts clientOptions) (*grpc.ClientConn, error) {

	uraryInts := []grpc.UnaryClientInterceptor{
		clientinterceptors.UnaryTimeoutInterceptor(opts.timeout), // 添加超时拦截器
	}
	if len(opts.unaryInterceptors) > 0 {
		uraryInts = append(uraryInts, opts.unaryInterceptors...) // 追加用户传入的拦截器
	}

	if opts.enableMetrics {
		uraryInts = append(uraryInts,
			clientinterceptors.UnaryPrometheusInterceptor(opts.histogramVecOpts, opts.counterVecOpts))
	}

	steamInts := []grpc.StreamClientInterceptor{}
	if len(opts.streamInterceptors) > 0 {
		steamInts = append(steamInts, opts.streamInterceptors...) // 追加用户传入的拦截器
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "` + opts.balancerName + `"}`),
		grpc.WithChainUnaryInterceptor(uraryInts...),
		grpc.WithChainStreamInterceptor(steamInts...),
	}
	if len(opts.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, opts.grpcOpts...) // 追加用户传入的选项
	}

	// 服务发现
	if opts.discovery != nil {
		grpcOpts = append(grpcOpts,
			grpc.WithResolvers(
				discovery.NewBuilder(
					opts.discovery,
					discovery.WithInsecure(opts.insecure))))
	} else {
		grpcOpts = append(grpcOpts,
			grpc.WithResolvers(direct.NewDirectBuilder()))
	}

	// 如果是不安全连接，添加不安全选项
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}

	if opts.enableTracing {
		grpcOpts = append(grpcOpts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	CliConn, err := grpc.NewClient(opts.endpoint, grpcOpts...)

	return CliConn, err
}
