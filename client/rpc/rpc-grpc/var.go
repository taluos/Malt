package grpc

import "time"

const (
	defaultClientName = " my rpc client"
	defautBalancer    = "round_robin"
	defaultAddress    = "127.0.0.1:8080"
	defaultTimeout    = 5 * time.Second
)
