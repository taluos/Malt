package serverinterceptors

import (
	"context"
	"runtime/debug"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/pkg/log"

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
	if ctx.Err() != nil {
		log.Errorf("Context error: %s, Panic: %+v\n\n%s", ctx.Err().Error(), r, debug.Stack())
	} else {
		log.Errorf("Panic occurred: %+v\n\n%s", r, debug.Stack())
	}
	// log.Errorf(ctx.Err().Error(), "%+v\n\n%s", r, debug.Stack())
	return errors.WithCode(code.ErrUnknow, "panic: %v", r)
}
