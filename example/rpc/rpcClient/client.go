package rpcclient

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rpcclient "github.com/taluos/Malt/server/rpc"
	grpcClient "github.com/taluos/Malt/server/rpc/rpc-grpc"

	pb "github.com/taluos/Malt/example/test_proto"
)

// Run 启动 gRPC 客户端并优雅关闭
func Run(ctx context.Context) error {
	// 创建 gRPC 客户端，可根据需要自定义连接地址、超时时间等
	c, err := rpcclient.NewClient("grpc",
		grpcClient.WithClientEndpoint("127.0.0.1:50051"),
		grpcClient.WithClientTimeout(5*time.Second),
		grpcClient.WithClientInsecure(true),
	)
	if err != nil {
		return err
	}
	// 创建 Greeter 客户端
	client := pb.NewGreeterClient(c.Conn().(*grpcClient.Client).ClientConn)

	// 调用 SayHello 方法
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Malt用户"})
	if err != nil {
		log.Printf("调用 SayHello 失败: %v", err)
	} else {
		log.Printf("收到服务器响应: %s", resp.GetMessage())
	}

	// 设置优雅关闭：监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭 gRPC 客户端...")

	// 优雅关闭客户端连接
	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		if err := c.Close(closeCtx); err != nil {
			log.Printf("关闭 gRPC 客户端出错: %v", err)
		}
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("gRPC 客户端已正常关闭")
	case <-time.After(5 * time.Second):
		log.Println("gRPC 客户端关闭超时，强制关闭")
	}

	return nil
}
