package httpserver

import (
	"os"

	maltAgent "github.com/taluos/Malt/core/trace"
	"github.com/taluos/Malt/pkg/errors"
	auth "github.com/taluos/Malt/server/rest/rest-gin/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type serverOptions struct {
	name    string `validate:"required"` // server name
	address string `validate:"required"` // server address:  IP_ADDRESS:PORT
	mode    string `validate:"required"` // gin.ReleaseMode, gin.DebugMode, gin.TestMode
	trans   string `validate:"required"` // http, https

	enableHealth    bool `validate:"required"` // healthz
	enableProfiling bool `validate:"required"` // pprof
	enableMetrics   bool `validate:"required"` // metrics
	enableTracing   bool `validate:"required"` // tracing
	enableCert      bool `validate:"required"` // https cert

	certFile string // https cert file
	keyFile  string // https key file

	trustedProxies []string          // trusted proxies
	middlewares    []gin.HandlerFunc // middlewares

	agent        *maltAgent.Agent   // tracing agent
	authOperator *auth.AuthOperator // auth operator
}

func (o *serverOptions) Validate() error {
	validator := validator.New()
	if err := validator.Struct(o); err != nil {
		return err
	}

	// 自定义验证逻辑
	if o.enableCert {
		if o.certFile == "" || o.keyFile == "" {
			return errors.New("cert file and key file are required when HTTPS is enabled")
		}
		// 验证文件是否存在
		if _, err := os.Stat(o.certFile); os.IsNotExist(err) {
			return errors.Errorf("cert file does not exist: %s", o.certFile)
		}
		if _, err := os.Stat(o.keyFile); os.IsNotExist(err) {
			return errors.Errorf("key file does not exist: %s", o.keyFile)
		}
	}

	return nil
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
