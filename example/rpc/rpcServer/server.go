package rpcserver

import (
	rpcserver "github.com/taluos/Malt/server/rpc/rpcServer"

	pb "github.com/taluos/Malt/example/test_proto"

	"github.com/taluos/Malt/example/rpc/service"

	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动 gRPC 服务并优雅关闭
func Run(ctx context.Context) error {
	// 创建 gRPC Server，可根据需要自定义监听地址、超时时间等
	s := rpcserver.NewServer(
		rpcserver.WithAddress("127.0.0.1:50051"),
		rpcserver.WithTimeout(5*time.Second),
	)

	pb.RegisterGreeterServer(s.Server, service.NewGreeterServer())

	// 启动服务（异步）
	go func() {
		log.Printf("gRPC server strating.")
		// 监听端口并启动服务
		if err := s.Start(ctx); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// 优雅关闭：监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")

	stopped := make(chan struct{})
	go func() {
		s.Stop(ctx)
		close(stopped)
	}()
	select {
	case <-stopped:
		log.Println("gRPC server exited properly")
	case <-time.After(5 * time.Second):
		log.Println("gRPC server shutdown timeout, forcing stop")
		s.Server.Stop()
	}
	return nil
}

func RPCInit() *rpcserver.Server {
	// 创建 gRPC Server，可根据需要自定义监听地址、超时时间等
	s := rpcserver.NewServer(
		rpcserver.WithAddress("127.0.0.1:50051"),
		rpcserver.WithTimeout(5*time.Second),
	)
	return s
}

func RPCRun(srv *rpcserver.Server, ctx context.Context) error {
	return srv.Start(ctx)
}

func RPCStop(srv *rpcserver.Server, ctx context.Context) error {
	<-ctx.Done() // 等待接收退出信号
	log.Println("Shutting down server...")
	sctx, cancal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancal()

	return srv.Stop(sctx)
}
