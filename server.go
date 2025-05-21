package Malt

import (
	"context"

	restServer "github.com/taluos/Malt/server/rest/Server"
	rpcServer "github.com/taluos/Malt/server/rpc/rpcServer"
)

var _ Server = (*rpcServer.Server)(nil)
var _ Server = (*restServer.Server)(nil)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
