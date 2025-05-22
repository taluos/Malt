package httpserver

import (
	"github.com/gin-gonic/gin"

	auth "github.com/taluos/Malt/core/auth"
	maltAgent "github.com/taluos/Malt/core/trace"
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

	trustedProxies []string
	middlewares    []gin.HandlerFunc

	// 添加认证操作器
	agent        *maltAgent.Agent
	authOperator *auth.AuthOperator
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

func WithTrustedProxies(trustedProxies []string) ServerOptions {
	return func(o *serverOptions) {
		o.trustedProxies = trustedProxies
	}
}

func WithMiddleware(middlewares ...gin.HandlerFunc) ServerOptions {
	return func(o *serverOptions) {
		o.middlewares = append(o.middlewares, middlewares...)
	}
}

func WithAgent(agent *maltAgent.Agent) ServerOptions {
	return func(o *serverOptions) {
		o.agent = agent
	}
}

func WithAuthOperator(authOperator *auth.AuthOperator) ServerOptions {
	return func(o *serverOptions) {
		o.authOperator = authOperator
	}
}
