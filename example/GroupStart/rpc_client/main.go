package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	consulApi "github.com/hashicorp/consul/api"
	rpcclient "github.com/taluos/Malt/client/rpc"
	grpcClient "github.com/taluos/Malt/client/rpc/rpc-grpc"
	consulRegistry "github.com/taluos/Malt/core/registry/consul"
	"github.com/taluos/Malt/core/selector"
	"github.com/taluos/Malt/core/selector/picker/random"
	pb "github.com/taluos/Malt/example/test_proto"
	"google.golang.org/grpc"
)

func main() {
	// 设置全局负载均衡策略
	selector.SetGlobalSelector(random.NewBuilder())
	grpcClient.InitBuilder()

	// 创建consul客户端
	conf := consulApi.DefaultConfig()
	conf.Address = "192.168.142.136:8500"
	conf.Scheme = "http"
	consulClient, err := consulApi.NewClient(conf)
	if err != nil {
		log.Fatalf("创建consul客户端失败: %v", err)
	}

	// 创建服务发现实例
	discovery := consulRegistry.New(
		consulClient,
		consulRegistry.WithHealthCheck(true),
		consulRegistry.WithHeartbeat(true),
		consulRegistry.WithHealthCheckInterval(10),
	)

	// 创建带服务发现的RPC客户端
	client, err := rpcclient.NewClient("grpc",
		grpcClient.WithEndpoint("discovery:///Malt-grpc"), // 使用服务发现
		grpcClient.WithBalancerName("random"),
		grpcClient.WithBalancerName("selector"),
		grpcClient.WithTimeout(5*time.Second),
		grpcClient.WithDiscovery(discovery), // 设置服务发现
	)
	if err != nil {
		log.Fatalf("创建RPC客户端失败: %v", err)
	}

	// 创建Greeter客户端
	conn, ok := client.Conn().(*grpc.ClientConn)
	if !ok {
		log.Fatalf("获取gRPC连接失败")
	}
	greeterClient := pb.NewGreeterClient(conn)

	// 创建上下文和取消函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动客户端调用
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 调用SayHello方法
				resp, err := greeterClient.SayHello(context.Background(), &pb.HelloRequest{
					Name: "RPC客户端用户",
				})
				if err != nil {
					log.Printf("调用SayHello失败: %v", err)
				} else {
					log.Printf("收到服务器响应: %s", resp.GetMessage())
				}
			}
		}
	}()

	log.Println("RPC客户端已启动，使用服务发现连接到Malt-grpc服务")
	log.Printf("客户端连接端点: %s", client.Endpoint())

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("正在关闭RPC客户端...")
	cancel()

	// 关闭客户端
	closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeCancel()
	if err := client.Close(closeCtx); err != nil {
		log.Printf("关闭RPC客户端失败: %v", err)
	} else {
		log.Println("RPC客户端已成功关闭")
	}
}
