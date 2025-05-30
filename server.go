package Malt

import (
	"context"

	maltServer "github.com/taluos/Malt/server"
)

var _ Server = (maltServer.Server)(nil)

type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error
	//	Stop 停止服务器
	Stop(ctx context.Context) error
}
