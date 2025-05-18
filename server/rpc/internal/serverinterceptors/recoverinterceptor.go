// This file is a copy of https://github.com/zeromicro/go-zero/blob/master/zrpc/internal/serverinterceptors/recoverinterceptor.go
// and modified by Malt.
// In detail : we change the log package used in func toPanicError.
package serverinterceptors

import (
	"context"
	"runtime/debug"

	"Malt/pkg/errors"
	"Malt/pkg/errors/code"
	"Malt/pkg/log"

	"google.golang.org/grpc"
)

// StreamRecoverInterceptor catches panics in processing stream requests and recovers.
func StreamRecoverInterceptor(svr any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r any) {
		err = toPanicError(context.Background(), r)
	})

	return handler(svr, stream)
}

// UnaryRecoverInterceptor catches panics in processing unary requests and recovers.
// this func is the old UnaryCrashInterception func.
// isolated the original clash handle func.
// in this func , we ignore the detail of defferent Server by ignore the info of grpc.UnaryServerInfo.
func UnaryRecoverInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer handleCrash(func(r any) {
		err = toPanicError(ctx, r)
	})

	return handler(ctx, req)
}

func handleCrash(handler func(any)) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(ctx context.Context, r any) error {
	log.Errorf(ctx.Err().Error(), "%+v\n\n%s", r, debug.Stack())
	return errors.WithCode(code.ErrUnknow, "panic: %v", r)
}
