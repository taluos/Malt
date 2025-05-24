package fiber

import (
	"context"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/log"
	"github.com/taluos/Malt/pkg/validations"
	middleware "github.com/taluos/Malt/server/rest/rest-fiber/internal/middlewares"
	"github.com/taluos/Malt/server/rest/rest-fiber/internal/pprof"

	uTranslator "github.com/go-playground/universal-translator"
	fiber "github.com/gofiber/fiber/v3"
)

// Wrapper for fiber.App
type Server struct {
	*fiber.App

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
		trans:   defaultrans,

		enableHealth:    true,
		enableProfiling: true,
		enableMetrics:   false,
		enableTracing:   false,

		trustedProxies: []string{},
		middlewares:    []fiber.Handler{},
	}

	// 先应用用户选项
	for _, opt := range opts {
		opt(o)
	}

	// 然后根据选项配置添加中间件
	if o.enableTracing && o.agent != nil {
		o.middlewares = append(o.middlewares, middleware.TracingMiddleware(o.agent))
	}

	if o.authOperator != nil {
		// 添加认证中间件
		o.middlewares = append(o.middlewares, middleware.AuthenticMiddleware(o.authOperator))
	}

	// 创建fiber配置
	config := fiber.Config{
		AppName: o.name,
	}

	// 创建服务器实例
	s := &Server{
		App:  fiber.New(config),
		opts: o,
	}

	// 应用中间件
	for _, mw := range o.middlewares {
		s.Use(mw)
	}

	// 配置健康检查
	if s.opts.enableHealth {
		s.Get("/health", func(c fiber.Ctx) error {
			return c.SendString("ok")
		})
	}

	// 配置性能分析
	if s.opts.enableProfiling {
		s.Use(pprof.PPofMiddleware())
	}

	// 配置指标监控
	if s.opts.enableMetrics {
		// TODO: 实现fiber的metrics中间件
		log.Infof("Metrics enabled for fiber server")
	}

	// 初始化翻译器
	var err error
	s.trans, err = initTrans(s.opts.trans)
	if err != nil {
		log.Errorf("init translator failed: %s", err.Error())
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

// start fiber server
func (s *Server) Start(ctx context.Context) error {
	log.Infof("[FIBER] server is running on %s", s.opts.address)

	s.rootCtx = ctx

	err := s.Listen(s.opts.address)
	if err != nil {
		return errors.Wrapf(err, "[FIBER] server failed")
	}

	return err
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("[FIBER] server is stopping on %s", s.opts.address)

	err := s.ShutdownWithContext(ctx)
	if err != nil {
		log.Errorf("[FIBER] server stopping failed: %s", err.Error())
		return errors.Wrapf(err, "[FIBER] stop server failed")
	}

	log.Infof("[FIBER] server is stopped on %s", s.opts.address)

	return err
}
