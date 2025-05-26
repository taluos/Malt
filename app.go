package Malt

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/taluos/Malt/core/registry"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/log"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
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
	ctx    context.Context
	cancel context.CancelFunc

	// 服务实例
	instances []*registry.ServiceInstance
	mu        sync.RWMutex

	opts options
}

func New(opts ...Option) *App {
	o := options{
		version:          Release,
		name:             defaultName,
		signal:           []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT},
		registrarTimeout: defalregistrarTimeout,
		stopTimeout:      defaltTimeout,
	}

	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	// 初始化上下文
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
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
	var endpoints []string
	for _, endpoint := range a.opts.endpoints {
		endpoints = append(endpoints, endpoint.String())
	}
	return endpoints
}

// 服务启动
func (app *App) Run() error {

	// 获取注册信息
	instances, err := app.buildInstance()
	if err != nil {
		return err
	}

	// 保护实列
	app.mu.Lock()
	app.instances = instances
	app.mu.Unlock()

	sctx := NewContext(app.ctx, app)
	eg, ctx := errgroup.WithContext(sctx)
	wg := sync.WaitGroup{}

	// start Rest Server
	for _, srv := range app.opts.restserver {
		server := srv
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx := sctx
			if app.opts.stopTimeout > 0 {
				var cancel context.CancelFunc
				stopCtx, cancel = context.WithTimeout(stopCtx, app.opts.stopTimeout)
				defer cancel()
			}
			return server.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return server.Start(sctx)
		})
	}

	// start RPC Server
	for _, srv := range app.opts.rpcserver {
		server := srv
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx := sctx
			if app.opts.stopTimeout > 0 {
				var cancel context.CancelFunc
				stopCtx, cancel = context.WithTimeout(stopCtx, app.opts.stopTimeout)
				defer cancel()
			}
			return server.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return server.Start(sctx)
		})
	}

	wg.Wait()

	// register service
	if app.opts.registrar != nil {
		rctx, cancel := context.WithTimeout(ctx, app.opts.registrarTimeout)
		defer cancel()
		for _, instance := range app.instances {
			err = app.opts.registrar.Register(rctx, instance)
			if err != nil {
				log.Errorf("register service error: %s", err)
				return err
			}
		}
	}

	// 监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, app.opts.signal...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return app.Stop()
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// 服务停止
func (app *App) Stop() error {
	app.mu.Lock()
	instances := app.instances
	app.mu.Unlock()

	if app.opts.registrar != nil && instances != nil {
		ctx, cancel := context.WithTimeout(NewContext(app.ctx, app), app.opts.stopTimeout)
		defer cancel()
		for _, instance := range app.instances {
			err := app.opts.registrar.Deregister(ctx, instance)
			if err != nil {
				log.Errorf("deregister service error: %s", err)
				return err
			}
		}

	}
	if app.cancel != nil {
		app.cancel()
	}

	return nil
}

// 创建服务注册结构体
func (app *App) buildInstance() ([]*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	//tags := make([]string, 0)
	instences := make([]*registry.ServiceInstance, 0)

	for _, endendpoint := range app.opts.endpoints {
		endpoints = append(endpoints, endendpoint.String())
		instences = append(instences, &registry.ServiceInstance{
			ID:        uuid.NewString(),
			Name:      app.opts.name + "-" + endendpoint.Scheme,
			Version:   app.opts.version,
			Endpoints: []string{endendpoint.String()},
			Metadata:  app.opts.metadata,
			Tags:      []string{endendpoint.Scheme},
		})
	}
	log.Infof("Create %d instances", len(instences))
	return instences, nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
