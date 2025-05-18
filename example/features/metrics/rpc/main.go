package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	rpcclient "github.com/taluos/Malt/example/features/metrics/rpc/client"
)

func main() {
	// 初始化客户端
	cli, err := rpcclient.RPCClientInit()
	if err != nil {
		log.Fatalf("初始化 RPC 客户端失败: %v", err)
	}

	// 使用客户端
	rpcclient.RPCClientUse(cli)

	// 设置信号处理，实现优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待退出信号或者主动退出
	select {
	case <-quit:
		log.Println("接收到退出信号，正在关闭客户端...")
	}

	// 关闭客户端
	if err := rpcclient.RPCClientClose(cli); err != nil {
		log.Fatalf("关闭 RPC 客户端失败: %v", err)
	}
}
