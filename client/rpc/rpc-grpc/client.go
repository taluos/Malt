package grpc

import (
	"context"

	interceptors "github.com/taluos/Malt/client/rpc/rpc-grpc/internal/interceptors"
	"github.com/taluos/Malt/core/resolver/direct"
	"github.com/taluos/Malt/core/resolver/discovery"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/log"

	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	// grpc.ClientConn
	*grpc.ClientConn

	rootCtx    context.Context
	rootCancel context.CancelFunc

	opts clientOptions
}

func NewClient(opts ...ClientOptions) (*Client, error) {
	var (
		err     error
		CliConn *grpc.ClientConn
	)

	o := clientOptions{
		name:    defaultClientName,
		address: defaultAddress,

		insecure:      true,
		enableTracing: false,
		enableMetrics: false,

		timeout:      defaultTimeout,
		balancerName: defautBalancer,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if err = o.Validate(); err != nil {
		log.Errorf("[gRPC] client options validate error: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	cli := &Client{
		rootCtx:    ctx,
		rootCancel: cancel,
		opts:       o,
	}

	CliConn, err = dial(o.insecure, o)
	if err != nil {
		cancel()
		return nil, err
	}
	cli.ClientConn = CliConn

	return cli, nil
}

func (c *Client) Close(ctx context.Context) error {

	if c.ClientConn == nil {
		return errors.New("client not initialized")
	}

	if c.rootCancel != nil {
		c.rootCancel()
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Infof("[gRPC] client closing")
		c.ClientConn.Close()
	}()

	select {
	case <-done:
	case <-ctx.Done():
		log.Errorf("[gRPC] server couldn't close gracefully in time")
		return ctx.Err()
	}

	log.Infof("[gRPC] client closed")
	return nil
}

func dial(insecure bool, opts clientOptions) (*grpc.ClientConn, error) {

	uraryInts := []grpc.UnaryClientInterceptor{
		interceptors.UnaryTimeoutInterceptor(opts.timeout), // 添加超时拦截器
	}
	if len(opts.unaryInterceptors) > 0 {
		uraryInts = append(uraryInts, opts.unaryInterceptors...) // 追加用户传入的拦截器
	}

	if opts.enableMetrics {
		uraryInts = append(uraryInts,
			interceptors.UnaryPrometheusInterceptor(opts.histogramVecOpts, opts.counterVecOpts))
	}

	if opts.enableTracing {
		//grpcOpts = append(grpcOpts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		uraryInts = append(uraryInts,
			interceptors.UnaryTracingInterceptor(opts.agent))
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
			grpc.WithResolvers(discovery.NewBuilder(
				opts.discovery,
				discovery.WithTimeout(opts.timeout),
				discovery.WithInsecure(insecure),
			),
			))
	} else {
		grpcOpts = append(grpcOpts,
			grpc.WithResolvers(direct.NewDirectBuilder()))
	}

	// 如果是不安全连接，添加不安全选项
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}

	CliConn, err := grpc.NewClient(opts.address, grpcOpts...)
	if err != nil {
		log.Errorf("[gRPC] dial error: %v", err)
		return nil, err
	}

	return CliConn, err
}

func (c *Client) Name() string {
	return c.opts.name
}

func (c *Client) Endpoint() string {
	return c.opts.address
}

func (c *Client) Balancer() string {
	return c.opts.balancerName
}
