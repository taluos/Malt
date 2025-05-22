package main

import (
	"context"
	"log"

	restserver "github.com/taluos/Malt/example/rest/restServer"
)

func main() {
	// easy start
	Server := restserver.RestInit()

	log.Println("Starting Rest server...")
	restserver.RestRun(Server, context.Background())
}
