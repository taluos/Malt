package rpcclient

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rpcclient "github.com/taluos/Malt/client/rpc"
	grpcClient "github.com/taluos/Malt/client/rpc/rpc-grpc"

	pb "github.com/taluos/Malt/example/test_proto"
	"google.golang.org/grpc"
)

// Run 启动 gRPC 客户端并优雅关闭
func Run(ctx context.Context) error {
	// 创建 gRPC 客户端，可根据需要自定义连接地址、超时时间等
	client, err := rpcclient.NewClient("grpc",
		grpcClient.WithEndpoint("127.0.0.1:50051"),
		grpcClient.WithTimeout(5*time.Second),
		grpcClient.WithInsecure(true),
	)
	if err != nil {
		return err
	}

	// 创建 Greeter 客户端
	conn := client.Conn().(*grpc.ClientConn)
	greeterClient := pb.NewGreeterClient(conn)

	// 调用 SayHello 方法
	resp, err := greeterClient.SayHello(ctx, &pb.HelloRequest{Name: "Malt用户"})
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

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err = client.Close(closeCtx); err != nil {
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

	return err
}
