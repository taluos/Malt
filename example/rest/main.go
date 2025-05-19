package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	restserver "github.com/taluos/Malt/example/rest/restServer"
)

func main() {
	// easy start

	Server := restserver.RestInit()

	restserver.InitRouter(Server)

	restserver.RestRun(Server, context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")

	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	restserver.RestStop(Server, sctx)

}
