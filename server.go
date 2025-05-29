package Malt

import (
	"context"

	maltServer "github.com/taluos/Malt/server"
	restServer "github.com/taluos/Malt/server/rest"
	rpcServer "github.com/taluos/Malt/server/rpc"
)

var _ maltServer.Server = (rpcServer.Server)(nil)
var _ maltServer.Server = (restServer.Server)(nil)

var _ Server = (maltServer.Server)(nil)

type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error
	//	Stop 停止服务器
	Stop(ctx context.Context) error
}
