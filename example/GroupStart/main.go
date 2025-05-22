package main

import (
	"time"

	consulApi "github.com/hashicorp/consul/api"
	malt "github.com/taluos/Malt"
	consulRegistry "github.com/taluos/Malt/core/registry/consul"
	"github.com/taluos/Malt/pkg/log"
	restserver "github.com/taluos/Malt/server/rest"
	ginServer "github.com/taluos/Malt/server/rest/rest-gin"
	rpcserver "github.com/taluos/Malt/server/rpc/rpcServer"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {

	restServerSet := []restserver.Server{}
	rpcServerSet := []rpcserver.Server{}

	restServerInstance := restserver.NewServer("gin",
		ginServer.WithPort(8080),
		ginServer.WithMiddleware(gin.Recovery()),
	)

	rpcServerInstance := rpcserver.NewServer(
		rpcserver.WithAddress("127.0.0.1:50051"),
		rpcserver.WithTimeout(5*time.Second),
	)

	restServerSet = append(restServerSet, restServerInstance)
	rpcServerSet = append(rpcServerSet, *rpcServerInstance)

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

	var App = malt.New(
		malt.WithId(uuid.New().String()),
		malt.WithName("Malt"),
		malt.WithTags([]string{"Rest:8080", "RPC:50051"}),
		malt.WithMetadata(map[string]string{"env": "dev", "Rest": "8080", "RPC": "50051"}),
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
