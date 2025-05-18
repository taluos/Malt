package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"Malt/pkg/errors"
	"Malt/pkg/log"
	"Malt/server/rest/internal/pprof"
	"Malt/server/rest/internal/validations"

	"github.com/gin-gonic/gin"
	uTranslator "github.com/go-playground/universal-translator"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Wrapper for gin.Engine
type Server struct {
	*gin.Engine
	server *http.Server
	trans  uTranslator.Translator

	opts *serverOptions
}

type ServerMethod interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func NewServer(opts ...ServerOptions) *Server {
	o := &serverOptions{
		name:            defaultName,
		port:            defaultPort,
		mode:            gin.DebugMode, // debug / release / test
		trans:           defaultrans,
		healthz:         true,
		enableProfiling: true,
		enableMetrics:   false,
		enableTracing:   false,
		middlewares:     []gin.HandlerFunc{},
	}

	o.jwt = NewJwtInfo()

	// 先应用用户选项
	for _, opt := range opts {
		opt(o)
	}

	// 然后根据选项配置添加中间件
	if o.enableTracing {
		o.middlewares = append(o.middlewares, otelgin.Middleware(o.name))
	}

	// 创建服务器实例
	s := &Server{
		Engine: gin.Default(),
		opts:   o,
	}

	// 应用中间件
	s.Use(o.middlewares...)

	// 配置健康检查
	if s.opts.healthz {
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
	// 设置gin模式
	gin.SetMode(s.opts.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	log.Infof("Rest server is running on port %d", s.opts.port)

	_ = s.SetTrustedProxies(nil)

	addr := fmt.Sprintf(":%d", s.opts.port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.Engine,
	}

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return errors.Wrapf(err, "start rest server failed")
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	var err error

	log.Infof("Rest server is stopping on port %d", s.opts.port)

	err = s.server.Shutdown(ctx)
	if err != nil {
		log.Errorf("stop rest server failed: %s", err.Error())
		return errors.Wrapf(err, "stop rest server failed")
	}

	log.Infof("Rest server is stopped on port %d", s.opts.port)

	return nil
}
