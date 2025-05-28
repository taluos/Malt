package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type (
	// MethodTimeoutConf defines specified timeout for gRPC method.
	MethodTimeoutConf struct {
		FullMethod string
		Timeout    time.Duration
	}
	methodTimeouts map[string]time.Duration
)

func UnaryTimeoutInterceptor(timeout time.Duration, methodTimeouts ...MethodTimeoutConf) grpc.UnaryClientInterceptor {
	timeouts := buildMethodTimeouts(methodTimeouts)
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取当前方法的超时时间
		t := getTimeoutByUnaryServerInfo(method, timeouts, timeout)
		if t <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		// 创建一个带有超时的上下文
		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel() // 确保在函数返回时取消上下文
		// 调用 invoker 函数，传递带有超时的上下文
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamTimeoutInterceptor(timeout time.Duration, methodTimeouts ...MethodTimeoutConf) grpc.StreamClientInterceptor {
	timeouts := buildMethodTimeouts(methodTimeouts)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// 获取当前方法的超时时间
		t := getTimeoutByUnaryServerInfo(method, timeouts, timeout)
		if t <= 0 {
			return streamer(ctx, desc, cc, method, opts...)
		}
		// 创建一个带有超时的上下文
		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel() // 确保在函数返回时取消上下文
		// 调用 streamer 函数，传递带有超时的上下文
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func buildMethodTimeouts(timeouts []MethodTimeoutConf) methodTimeouts {
	// 构建方法超时配置映射
	mt := make(methodTimeouts, len(timeouts))
	for _, st := range timeouts {
		if st.FullMethod != "" {
			mt[st.FullMethod] = st.Timeout
		}
	}

	return mt
}

func getTimeoutByUnaryServerInfo(method string, timeouts methodTimeouts, defaultTimeout time.Duration) time.Duration {
	// 如果找到特定方法的超时配置，则返回该配置，否则返回默认超时时间
	if v, ok := timeouts[method]; ok {
		return v
	}
	return defaultTimeout
}
