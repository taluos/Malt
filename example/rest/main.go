package main

import (
	"context"
	"fmt"
	"log"
	"time"

	restclient "github.com/taluos/Malt/example/rest/client"
	restserver "github.com/taluos/Malt/example/rest/ginServer"
	//restserver "github.com/taluos/Malt/example/rest/fiberServer"
)

func main() {
	// easy start
	Server := restserver.RestInit()
	Client := restclient.ClientInit()

	go func() {
		log.Println("Starting Rest server...")
		restserver.RestRun(Server, context.Background())
	}()

	time.Sleep(5 * time.Second)

	responses, err := restclient.ClientGet(Client, context.Background(), "hello")
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	fmt.Printf("响应: %v", responses.String())

	time.Sleep(5 * time.Second)
}
