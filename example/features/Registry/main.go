package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	malt "github.com/taluos/Malt"
	consulRegistry "github.com/taluos/Malt/core/registry/consul"
	restserver "github.com/taluos/Malt/example/features/Registry/restServer"

	"github.com/hashicorp/consul/api"
)

const testConsulAddr = "192.168.142.136:8500"

func main() {

	Server := restserver.RestInit()
	restserver.InitRouter(Server)
	consulClient, err := api.NewClient(&api.Config{Address: testConsulAddr})
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
		malt.WithRegistrarTimeout(5*time.Second),
		malt.WithStopTimeout(5*time.Second),

		malt.WithRESTServer(*Server),
		malt.WithRegistrar(RegistyInstance),
	)

	err = App.Run()
	if err != nil {
		log.Fatalf("server failed: %v", err)
		panic(err)
	}
}
