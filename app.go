package Malt

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/taluos/Malt/core/registry"
	"github.com/taluos/Malt/pkg/log"

	"github.com/google/uuid"
)

// AppInfo is application context value.
type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}
type App struct {
	opts   options
	ctx    context.Context
	cancel context.CancelFunc

	// 服务实例
	instance *registry.ServiceInstance

	mu sync.RWMutex
}

func New(opts ...Option) *App {
	o := options{
		name:             defaultName,
		signal:           []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT},
		registrarTimeout: defalregistrarTimeout,
		stopTimeout:      defaltTimeout,
	}

	if o.id == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			log.Errorf("generate uuid error: %s", err)
			return nil
		}
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &App{
		opts: o,
	}
}

func (a *App) ID() string {
	return a.opts.id
}

func (a *App) Name() string {
	return a.opts.name
}

func (a *App) Version() string {
	return a.opts.version
}

func (a *App) Metadata() map[string]string {
	return a.opts.metadata
}

func (a *App) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

// 服务启动
func (app *App) Run() error {
	// 获取注册信息
	instance, err := app.buildInstance()
	if err != nil {
		return err
	}

	// 保护实列
	app.mu.Lock()
	app.instance = instance
	app.mu.Unlock()

	// register service
	if app.opts.registrar != nil {
		rctx, cancel := context.WithTimeout(context.Background(), app.opts.registrarTimeout)
		defer cancel()
		err = app.opts.registrar.Register(rctx, instance)
		if err != nil {
			log.Errorf("register service error: %s", err)
			return err
		}
	}

	// 监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, app.opts.signal...)
	<-c

	return nil
}

// 服务停止
func (app *App) Stop() error {
	app.mu.Lock()
	instance := app.instance
	app.mu.Unlock()

	if app.opts.registrar != nil && app.instance != nil {
		rctx, cancel := context.WithTimeout(context.Background(), app.opts.stopTimeout)
		defer cancel()
		err := app.opts.registrar.Deregister(rctx, instance)
		if err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}

	if app.cancel != nil {
		app.cancel()
	}

	return nil
}

// 创建服务注册结构体
func (app *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	tags := make([]string, 0)

	for _, endendpoint := range app.opts.endpoints {
		endpoints = append(endpoints, endendpoint.String())
	}

	tags = append(tags, app.opts.tags...)

	instance := &registry.ServiceInstance{
		ID:        app.opts.id,
		Name:      app.opts.name,
		Endpoints: endpoints,
		Tags:      tags,
		Version:   app.opts.version,
	}

	return instance, nil
}
