// this file is a modified of https://github.com/zeromicro/go-zero/blob/master/zrpc/internal/serverinterceptors/timeoutinterceptor.go
// we modify the error handling
package serverinterceptors

import (
	"github.com/taluos/Malt/pkg/errors"

	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// MethodTimeoutConf defines specified timeout for gRPC method.
	MethodTimeoutConf struct {
		FullMethod string
		Timeout    time.Duration
	}

	methodTimeouts map[string]time.Duration
)

// UnaryTimeoutInterceptor returns a func that sets timeout to incoming unary requests.
// Use closure to transfer the methodTimeouts to the inner func.
func UnaryTimeoutInterceptor(timeout time.Duration, methodTimeouts ...MethodTimeoutConf) grpc.UnaryServerInterceptor {
	timeouts := buildMethodTimeouts(methodTimeouts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 获取当前方法的超时时间
		t := getTimeoutByUnaryServerInfo(info.FullMethod, timeouts, timeout)
		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		var resp any
		var err error
		var lock sync.Mutex
		done := make(chan struct{})
		// 创建带缓冲的panic通道，避免goroutine泄漏
		panicChan := make(chan any, 1)

		// 在新的goroutine中执行handler
		// output the result to the done and panicChan
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// 附加调用栈信息，避免在不同goroutine中丢失
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack()))) // 将错误放到panicChan中
				}
			}()
			// lock the handler to avoid race condition
			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done)
		}()

		// 在 goroutine 外对执行结果做处理
		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			// Same as the handler, have to lock the resp to avoid race condition
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done():
			// if handler is out of time, return the error
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				// if the error is canceled, return the canceled error
				// mostly because the ctx is canceled by other interceptor
				// err = errors.WithCode(err, codes.Canceled)
				err = status.Error(codes.Canceled, err.Error())
			} else if errors.Is(err, context.DeadlineExceeded) {
				// if the error is deadline exceeded, return the deadline exceeded error
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, err
		}
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
