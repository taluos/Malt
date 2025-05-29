package client

import "time"

const (
	RESTClientType = "rest"
	RPCClientType  = "rpc"
)

var (
	defaultTimeout = 5 * time.Second
)
