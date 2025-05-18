package client

import (
	"time"

	rpcclient "github.com/taluos/Malt/server/rpc/rpcClient"

	pb "github.com/taluos/Malt/example/test_proto"

	"context"
	"log"
)

// RPCClientInit 初始化 RPC 客户端
func RPCClientInit() (*rpcclient.Client, error) {
	// 创建 gRPC 客户端，可根据需要自定义连接地址、超时时间等
	c, err := rpcclient.NewClient(
		rpcclient.WithEndpoint("127.0.0.1:8090"),
		rpcclient.WithTimeout(5*time.Second),
		rpcclient.WithInsecure(true),
		rpcclient.WithEnableMetrics(true),
	)
	return c, err
}

// RPCClientClose 关闭 RPC 客户端
func RPCClientClose(cli *rpcclient.Client) error {
	log.Println("关闭 RPC 客户端...")
	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return cli.Close(closeCtx)
}

// RPCClientUse 使用 RPC 客户端示例
func RPCClientUse(cli *rpcclient.Client) {
	// 创建 Greeter 客户端
	client := pb.NewGreeterClient(cli.CliConn)

	// 调用 SayHello 方法
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Malt用户"})
	if err != nil {
		log.Printf("调用 SayHello 失败: %v", err)
	} else {
		log.Printf("收到服务器响应: %s", resp.GetMessage())
	}

	log.Println("使用 RPC 客户端连接:", cli.Endpoint())
}
