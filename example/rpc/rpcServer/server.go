package rpcserver

import (
	rpcserver "github.com/taluos/Malt/server/rpc"
	grpcServer "github.com/taluos/Malt/server/rpc/rpc-grpc"

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
	s := rpcserver.NewServer("grpc",
		grpcServer.WithAddress("127.0.0.1:50051"),
		grpcServer.WithTimeout(5*time.Second),
	)

	s.RegisterService(pb.RegisterGreeterServer, service.NewGreeterServer())

	//pb.RegisterGreeterServer(s.Server, service.NewGreeterServer())

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

	stopped := make(chan struct{}, 1)
	go func() {
		s.Stop(ctx)
		close(stopped)
	}()
	select {
	case <-stopped:
		log.Println("gRPC server exited properly")
	case <-time.After(5 * time.Second):
		log.Println("gRPC server shutdown timeout, forcing stop")
		s.Stop(ctx)
		// s.Server.Stop()
	}
	return nil
}
