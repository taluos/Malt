package fiber

import (
	fiber "github.com/gofiber/fiber/v3"

	maltAgent "github.com/taluos/Malt/core/trace"
	auth "github.com/taluos/Malt/server/rest/rest-fiber/internal/auth"
)

type serverOptions struct {
	name    string
	address string

	trans string

	enableHealth    bool
	enableProfiling bool
	enableMetrics   bool
	enableTracing   bool

	trustedProxies []string
	middlewares    []fiber.Handler

	agent        *maltAgent.Agent
	authOperator *auth.AuthOperator
}

type ServerOptions func(*serverOptions)

func WithName(name string) ServerOptions {
	return func(o *serverOptions) {
		o.name = name
	}
}

func WithAddress(address string) ServerOptions {
	return func(o *serverOptions) {
		o.address = address
	}
}

func WithTrans(trans string) ServerOptions {
	return func(o *serverOptions) {
		o.trans = trans
	}
}

func WithHealthz(healthz bool) ServerOptions {
	return func(o *serverOptions) {
		o.enableHealth = healthz
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

func WithMiddleware(middlewares ...fiber.Handler) ServerOptions {
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
