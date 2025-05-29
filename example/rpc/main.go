package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rpcserver "github.com/taluos/Malt/example/rpc/rpcServer"

	rpcclient "github.com/taluos/Malt/example/rpc/rpcClient"
)

func main() {
	// 创建上下文
	ctx := context.Background()

	// 启动服务器（异步）
	go func() {
		if err := rpcserver.Run(ctx); err != nil {
			log.Fatalf("服务器运行失败: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(2 * time.Second)

	// 启动客户端
	go func() {
		if err := rpcclient.Run(ctx); err != nil {
			log.Fatalf("客户端运行失败: %v", err)
		}
	}()

	// 阻塞主线程，等待客户端和服务器完成
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case <-stop:
		log.Println("应用已正常关闭")
	case <-time.After(5 * time.Second):
		log.Println("应用关闭超时，强制关闭")
	}
}
