package httpserver

import (
	"context"
	"net/http"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/log"
	"github.com/taluos/Malt/pkg/validations"
	middleware "github.com/taluos/Malt/server/rest/rest-gin/internal/middlewares"
	"github.com/taluos/Malt/server/rest/rest-gin/internal/pprof"

	"github.com/gin-gonic/gin"
	uTranslator "github.com/go-playground/universal-translator"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

// Wrapper for gin.Engine
type Server struct {
	*gin.Engine

	server  *http.Server
	rootCtx context.Context
	trans   uTranslator.Translator

	opts *serverOptions
}

type ServerMethod interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func NewServer(opts ...ServerOptions) *Server {
	o := &serverOptions{
		name:    defaultName,
		address: defaultAddr,
		mode:    gin.DebugMode, // debug / release / test
		trans:   defaultrans,

		enableHealth:    true,
		enableProfiling: true,
		enableMetrics:   false,
		enableTracing:   false,
		enableCert:      false,

		certFile: "",
		keyFile:  "",

		trustedProxies: []string{},
		middlewares:    []gin.HandlerFunc{},
	}

	// 先应用用户选项
	for _, opt := range opts {
		opt(o)
	}

	// 然后根据选项配置添加中间件
	if o.enableTracing && o.agent != nil {
		//o.middlewares = append(o.middlewares, otelgin.Middleware(o.name))
		o.middlewares = append(o.middlewares, middleware.TracingMiddleware(o.agent))
	}

	if o.authOperator != nil {
		// 添加认证中间件
		o.middlewares = append(o.middlewares, middleware.AuthenticMiddleware(o.authOperator))
	}

	// 创建服务器实例
	s := &Server{
		Engine: gin.Default(),
		opts:   o,
	}

	// 应用中间件
	s.Use(o.middlewares...)

	// 配置健康检查
	if s.opts.enableHealth {
		s.GET("/health", func(c *gin.Context) {
			c.String(200, "ok")
		})
	}

	// 配置性能分析
	if s.opts.enableProfiling {
		pprof.Register(s.Engine)
	}

	// 配置指标监控
	if s.opts.enableMetrics {
		m := ginmetrics.GetMonitor()
		m.SetMetricPath("/metrics")
		m.SetSlowTime(5)
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		m.Use(s)
	}

	// 初始化翻译器
	var err error
	s.trans, err = initTrans(s.opts.trans)
	if err != nil {
		log.Errorf("init translator failed: %s", err.Error())
		// 这里可以考虑返回错误而不是继续
	}

	// 注册验证器
	if s.trans != nil {
		validations.RegisterEmail(s.trans)
		validations.RegisterMobile(s.trans)
		validations.RegisterPassword(s.trans)
		validations.RegisterUsername(s.trans)
	}

	return s
}

// statrt rest server
func (s *Server) Start(ctx context.Context) error {
	var err error
	// 设置gin模式
	gin.SetMode(s.opts.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	log.Infof("[HTTP] server is running on %v", s.opts.address)

	_ = s.SetTrustedProxies(s.opts.trustedProxies)

	// addr := fmt.Sprintf(":%d", s.opts.address)

	s.rootCtx = ctx
	s.server = &http.Server{
		Addr:    s.opts.address,
		Handler: s.Engine,
	}

	if s.opts.enableCert && s.opts.certFile != "" && s.opts.keyFile != "" {
		err = s.server.ListenAndServeTLS(s.opts.certFile, s.opts.keyFile)
	} else {
		err = s.server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		return errors.Wrapf(err, "[HTTP] server failed")
	}

	return err
}

func (s *Server) Stop(ctx context.Context) error {
	var err error

	log.Infof("[HTTP] server is stopping on %v", s.opts.address)

	err = s.server.Shutdown(ctx)
	if err != nil {
		log.Errorf("[HTTP] server stopping failed: %s", err.Error())
		return errors.Wrapf(err, "[HTTP] stop server failed")
	}

	log.Infof("[HTTP] server is stopped on %v", s.opts.address)

	return err
}
