package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/taluos/Malt/example/features/metrics/rpc/service"
	pb "github.com/taluos/Malt/example/test_proto"
	"github.com/taluos/Malt/pkg/log"
	rpcserver "github.com/taluos/Malt/server/rpc"
	grpcServer "github.com/taluos/Malt/server/rpc/rpc-grpc"
)

func rpcServerInit() rpcserver.Server {
	// 创建 gRPC Server，可根据需要自定义监听地址、超时时间等
	s := rpcserver.NewServer("grpc",
		grpcServer.WithAddress("127.0.0.1:8090"),
		grpcServer.WithTimeout(5*time.Second),
		grpcServer.WithEnableMetrics(true),
	)

	// 注册服务
	s = s.RegisterService(pb.RegisterGreeterServer, &service.GreeterServer{})
	// pb.RegisterGreeterServer(s.Server, &service.GreeterServer{})

	return s
}

func rpcRun(srv rpcserver.Server, ctx context.Context) error {
	return srv.Start(ctx)
}

func rpcStop(srv rpcserver.Server, ctx context.Context) error {
	<-ctx.Done()
	log.Info("RPC server stopping")
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Stop(sctx)
}

func main() {
	rpcServer := rpcServerInit()

	// 创建上下文和取消函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器（异步）
	go func() {
		log.Info("RPC server starting")
		if err := rpcRun(rpcServer, ctx); err != nil {
			log.Fatalf("RPC server stopped with error: %v", err)
		}
	}()

	// 等待退出信号
	<-quit
	log.Info("Received shutdown signal")

	// 取消上下文，触发服务停止
	cancel()

	// 优雅关闭服务器
	if err := rpcStop(rpcServer, context.Background()); err != nil {
		log.Errorf("Failed to stop RPC server: %v", err)
	}
}
