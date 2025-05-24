package serverinterceptors

import (
	"context"
	"time"

	metric "github.com/taluos/Malt/core/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	// ServerNamespace defines a logical grouping of servers.
	ServerNamespace = "rpc_server"
)

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: ServerNamespace,
		Subsystem: "requests",
		Name:      "requests_duration_seconds",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: ServerNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc server requests code count.",
		Labels:    []string{"method", "code"},
	})
)

func UnaryPrometheusInterceptor(histogramVecOpts *metric.HistogramVecOpts, counterVecOpts *metric.CounterVecOpts) grpc.UnaryServerInterceptor {

	if histogramVecOpts != nil {
		metricServerReqDur = metric.NewHistogramVec(histogramVecOpts)
	}

	if counterVecOpts != nil {
		metricServerReqCodeTotal = metric.NewCounterVec(counterVecOpts)
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		now := time.Now()

		resp, err := handler(ctx, req)

		// 记录耗时
		metricServerReqDur.Observe(int64(time.Since(now)/time.Millisecond), info.FullMethod)

		// 记录状态码
		metricServerReqCodeTotal.Inc(info.FullMethod, status.Code(err).String())

		return resp, err
	}
}
