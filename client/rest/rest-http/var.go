package http

import "time"

var (
	defaultTimeout = 30 * time.Second
	defaultCount   = 3
	defaultAgent   = "Malt-HTTP-Client"
)
