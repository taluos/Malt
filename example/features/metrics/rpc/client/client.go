package client

import (
	"context"
	"log"
	"time"

	rpcclient "github.com/taluos/Malt/client/rpc"
	grpcClient "github.com/taluos/Malt/client/rpc/rpc-grpc"
	pb "github.com/taluos/Malt/example/test_proto"

	"google.golang.org/grpc"
)

// RPCClientInit 初始化 RPC 客户端
func RPCClientInit() (rpcclient.Client, error) {
	// 创建 gRPC 客户端，可根据需要自定义连接地址、超时时间等
	c, err := rpcclient.NewClient("grpc",
		grpcClient.WithEndpoint("127.0.0.1:8090"),
		grpcClient.WithTimeout(5*time.Second),
		grpcClient.WithInsecure(true),
		grpcClient.WithEnableMetrics(true),
	)
	return c, err
}

// RPCClientClose 关闭 RPC 客户端
func RPCClientClose(cli rpcclient.Client) error {
	log.Println("关闭 RPC 客户端...")
	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return cli.Close(closeCtx)
}

// RPCClientUse 使用 RPC 客户端示例
func RPCClientUse(cli rpcclient.Client) {
	// 创建 Greeter 客户端
	// 创建 Greeter 客户端
	conn, _ := cli.Conn().(*grpc.ClientConn)
	client := pb.NewGreeterClient(conn)

	// 调用 SayHello 方法
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Malt用户"})
	if err != nil {
		log.Printf("调用 SayHello 失败: %v", err)
	} else {
		log.Printf("收到服务器响应: %s", resp.GetMessage())
	}

	log.Println("使用 RPC 客户端连接:", cli.Endpoint())
}
