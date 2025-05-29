package server

import (
	"context"
	"fmt"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/server/rest"
)

type restServerWrapper struct {
	server     rest.Server
	serverType string
}

var _ Server = (*restServerWrapper)(nil)

// NewRESTServer 创建REST服务器
func newRestServer(serverType string, opts ...rest.ServerOptions) (Server, error) {
	server := rest.NewServer(serverType, opts...)
	if server == nil {
		return nil, errors.New(fmt.Sprintf("[Server] failed to create REST server with type: %s", serverType))
	}
	return &restServerWrapper{server: server, serverType: serverType}, nil
}

func (w *restServerWrapper) Type() string {
	return w.serverType
}

func (w *restServerWrapper) Start(ctx context.Context) error {
	return w.server.Start(ctx)
}

func (w *restServerWrapper) Stop(ctx context.Context) error {
	return w.server.Stop(ctx)
}

// 嵌入REST服务器的所有方法
func (w *restServerWrapper) Group(relativePath string, handlers ...any) rest.RouteGroup {
	return w.server.Group(relativePath, handlers...)
}

func (w *restServerWrapper) Use(middleware ...any) rest.Server {
	return w.server.Use(middleware...)
}

func (w *restServerWrapper) Handle(httpMethod, relativePath string, handlers ...any) rest.Server {
	return w.server.Handle(httpMethod, relativePath, handlers...)
}
