package clientinterceptors

import (
	"context"
	"time"

	metric "github.com/taluos/Malt/core/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	// ServerNamespace defines a logical grouping of servers.
	ServerNamespace = "rpc_client"
)

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: ServerNamespace,
		Subsystem: "requests",
		Name:      "requests_duration_seconds",
		Help:      "rpc client requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: ServerNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc client requests code count.",
		Labels:    []string{"method", "code"},
	})
)

func UnaryPrometheusInterceptor(histogramVecOpts *metric.HistogramVecOpts, counterVecOpts *metric.CounterVecOpts) grpc.UnaryClientInterceptor {
	if histogramVecOpts != nil {
		metricServerReqDur = metric.NewHistogramVec(histogramVecOpts)
	}
	if counterVecOpts != nil {
		metricServerReqCodeTotal = metric.NewCounterVec(counterVecOpts)
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		// 记录耗时
		metricServerReqDur.Observe(int64(time.Since(now)/time.Millisecond), method)
		// 记录状态码
		metricServerReqCodeTotal.Inc(method, status.Code(err).String())
		return err
	}
}

func StreamPrometheusInterceptor(histogramVecOpts *metric.HistogramVecOpts, counterVecOpts *metric.CounterVecOpts) grpc.StreamClientInterceptor {
	if histogramVecOpts != nil {
		metricServerReqDur = metric.NewHistogramVec(histogramVecOpts)
	}
	if counterVecOpts != nil {
		metricServerReqCodeTotal = metric.NewCounterVec(counterVecOpts)
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		now := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		// 记录耗时
		metricServerReqDur.Observe(int64(time.Since(now)/time.Millisecond), method)
		// 记录状态码
		metricServerReqCodeTotal.Inc(method, status.Code(err).String())
		return clientStream, err
	}
}
