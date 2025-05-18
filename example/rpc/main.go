package main

import (
	"context"
	"log"
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
	select {}
}
