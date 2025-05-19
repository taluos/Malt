package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	malt "github.com/taluos/Malt"
	restserver "github.com/taluos/Malt/example/rest/restServer"
	rpcserver "github.com/taluos/Malt/example/rpc/rpcServer"
	"github.com/taluos/Malt/pkg/log"

	"golang.org/x/sync/errgroup"
)

func main() {

	restServer := restserver.RestInit()
	rpcServer := rpcserver.RPCInit()

	var _ malt.Server = restServer
	var _ malt.Server = rpcServer

	var servers = []malt.Server{}

	servers = append(servers, restServer)
	servers = append(servers, rpcServer)

	ctx, cancel := context.WithCancel(context.Background())

	eg, ctx := errgroup.WithContext(ctx)

	for _, srv := range servers {

		eg.Go(func() error {
			<-ctx.Done()
			log.Infof("Shutting down server...")
			sctx, cancal := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancal()
			return srv.Stop(sctx)
		})

		eg.Go(func() error {
			return srv.Start(ctx)
		})
	}

	// 优雅关闭：监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-quit:
			cancel()
			return nil
		case <-time.After(10 * time.Second):
			log.Fatal("timeout")
			return ctx.Err()
		}
	})

	if err := eg.Wait(); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
