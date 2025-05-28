package main

import (
	"fmt"
	"net/url"
	"time"

	consulApi "github.com/hashicorp/consul/api"
	malt "github.com/taluos/Malt"
	consulRegistry "github.com/taluos/Malt/core/registry/consul"
	"github.com/taluos/Malt/example/GroupStart/service"
	pb "github.com/taluos/Malt/example/test_proto"
	"github.com/taluos/Malt/pkg/log"
	getport "github.com/taluos/Malt/pkg/port"
	restserver "github.com/taluos/Malt/server/rest"
	ginServer "github.com/taluos/Malt/server/rest/rest-gin"
	rpcserver "github.com/taluos/Malt/server/rpc"
	grpcServer "github.com/taluos/Malt/server/rpc/rpc-grpc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {

	restServerSet := []restserver.Server{}
	rpcServerSet := []rpcserver.Server{}

	port1, _ := getport.GetFreePort()
	address1 := fmt.Sprintf("%s:%d", "10.60.82.146", port1)
	restServerInstance := restserver.NewServer("gin",
		ginServer.WithHealthz(true),
		ginServer.WithAddress(address1),
		ginServer.WithMiddleware(gin.Recovery()),
	)

	port2, _ := getport.GetFreePort()
	address2 := fmt.Sprintf("%s:%d", "10.60.82.146", port2)
	rpcServerInstance := rpcserver.NewServer("grpc",
		grpcServer.WithAddress(address2),
		grpcServer.WithTimeout(5*time.Second),
	)
	rpcServerInstance.RegisterService(pb.RegisterGreeterServer, service.NewGreeterServer())

	restServerSet = append(restServerSet, restServerInstance)
	rpcServerSet = append(rpcServerSet, rpcServerInstance)

	consulClient, err := consulApi.NewClient(&consulApi.Config{Address: "192.168.142.136:8500"})
	if err != nil {
		log.Fatalf("创建consul客户端失败: %v", err)
	}

	RegistyInstance := consulRegistry.New(
		consulClient,
		consulRegistry.WithHealthCheck(true),
		consulRegistry.WithHeartbeat(true),
		consulRegistry.WithHealthCheckInterval(10),
	)

	// 创建URL端点
	restEndpoint, _ := url.Parse(fmt.Sprintf("http://%s", address1))
	// restEndpoint.Query().Add("insecure", "true")
	rpcEndpoint, _ := url.Parse(fmt.Sprintf("grpc://%s", address2))
	// rpcEndpoint.Query().Add("insecure", "true")

	var App = malt.New(
		malt.WithId(uuid.New().String()),
		malt.WithName("Malt"),
		// 添加Endpoints
		malt.WithEndpoints([]*url.URL{restEndpoint, rpcEndpoint}),
		// 可以保留元数据
		malt.WithRegistrarTimeout(5*time.Second),
		malt.WithStopTimeout(5*time.Second),

		malt.WithRESTServer(restServerSet...),
		malt.WithRPCServer(rpcServerSet...),
		malt.WithRegistrar(RegistyInstance),
	)

	err = App.Run()
	if err != nil {
		log.Fatalf("server failed: %v", err)
		panic(err)
	}

}
