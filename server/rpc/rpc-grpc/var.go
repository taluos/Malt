package grpc

import "time"

const (
	defaultServerName = " my rpc server"
	defautBalancer    = "round_robin"
	defaultAddress    = "127.0.0.1:8080"
	defaultTimeout    = 5 * time.Second
)
