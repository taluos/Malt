package client

import (
	"time"

	pb "github.com/taluos/Malt/example/test_proto"
	rpcclient "github.com/taluos/Malt/server/rpc"
	grpcClient "github.com/taluos/Malt/server/rpc/rpc-grpc"
	"google.golang.org/grpc"

	"context"
	"log"
)

// RPCClientInit 初始化 RPC 客户端
func RPCClientInit() (rpcclient.Client, error) {
	// 创建 gRPC 客户端，可根据需要自定义连接地址、超时时间等
	c, err := rpcclient.NewClient("grpc",
		grpcClient.WithClientEndpoint("127.0.0.1:8090"),
		grpcClient.WithClientTimeout(5*time.Second),
		grpcClient.WithClientInsecure(true),
		grpcClient.WithClientEnableMetrics(true),
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
