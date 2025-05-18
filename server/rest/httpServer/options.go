package httpserver

import (
	"github.com/gin-gonic/gin"
)

type serverOptions struct {
	name string
	port int

	mode  string
	trans string

	healthz         bool
	enableProfiling bool
	enableMetrics   bool
	enableTracing   bool

	middlewares []gin.HandlerFunc
	jwt         *JwtInfo
}

type ServerOptions func(*serverOptions)

func WithName(name string) ServerOptions {
	return func(o *serverOptions) {
		o.name = name
	}
}

func WithPort(port int) ServerOptions {
	return func(o *serverOptions) {
		o.port = port
	}
}

func WithMode(mode string) ServerOptions {
	return func(o *serverOptions) {
		o.mode = mode
	}
}

func WithTrans(trans string) ServerOptions {
	return func(o *serverOptions) {
		o.trans = trans
	}
}

func WithHealthz(healthz bool) ServerOptions {
	return func(o *serverOptions) {
		o.healthz = healthz
	}
}

func WithEnableProfiling(enableProfiling bool) ServerOptions {
	return func(o *serverOptions) {
		o.enableProfiling = enableProfiling
	}
}

func WithEnableMetrics(enableMetrics bool) ServerOptions {
	return func(o *serverOptions) {
		o.enableMetrics = enableMetrics
	}
}

func WithEnableTracing(enableTracing bool) ServerOptions {
	return func(o *serverOptions) {
		o.enableTracing = enableTracing
	}
}

func WithMiddleware(middlewares ...gin.HandlerFunc) ServerOptions {
	return func(o *serverOptions) {
		o.middlewares = append(o.middlewares, middlewares...)
	}
}

func WithJwt(jwt *JwtInfo) ServerOptions {
	return func(o *serverOptions) {
		o.jwt = jwt
	}
}
