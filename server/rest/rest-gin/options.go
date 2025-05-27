package httpserver

import (
	"github.com/gin-gonic/gin"

	maltAgent "github.com/taluos/Malt/core/trace"
	auth "github.com/taluos/Malt/server/rest/rest-gin/internal/auth"
)

type serverOptions struct {
	name    string // server name
	address string // server address:  IP_ADDRESS:PORT

	mode  string // gin.ReleaseMode, gin.DebugMode, gin.TestMode
	trans string // http, https

	enableHealth    bool // healthz
	enableProfiling bool // pprof
	enableMetrics   bool // metrics
	enableTracing   bool // tracing
	enableCert      bool // https cert

	certFile string // https cert file
	keyFile  string // https key file

	trustedProxies []string          // trusted proxies
	middlewares    []gin.HandlerFunc // middlewares

	agent        *maltAgent.Agent   // tracing agent
	authOperator *auth.AuthOperator // auth operator
}

// ServerOptions is a function that takes a pointer to a serverOptions struct and modifies it.
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

func WithEnableCert(enableCert bool) ServerOptions {
	return func(o *serverOptions) {
		o.enableCert = enableCert
	}
}

func WithCertFile(certFile string) ServerOptions {
	return func(o *serverOptions) {
		o.certFile = certFile
	}
}

func WithKeyFile(keyFile string) ServerOptions {
	return func(o *serverOptions) {
		o.keyFile = keyFile
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
