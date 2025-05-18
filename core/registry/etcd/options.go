package etcd

import (
	"context"
	"time"
)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
	maxRetry  int
}

// Option is etcd registry option.
type Option func(o *options)

// Context with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Namespace with registry namespace.
func WithNamespace(ns string) Option {
	return func(o *options) { o.namespace = ns }
}

// RegisterTTL with register ttl.
func WithRegisterTTL(ttl time.Duration) Option {
	return func(o *options) { o.ttl = ttl }
}

func WithMaxRetry(num int) Option {
	return func(o *options) { o.maxRetry = num }
}
