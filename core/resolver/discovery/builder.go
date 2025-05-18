package discovery

import (
	"Malt/core/registry"
	"Malt/pkg/errors"

	"context"
	"strings"
	"time"

	"google.golang.org/grpc/resolver"
)

var _ resolver.Builder = (*builder)(nil)

type builder struct {
	discovery registry.Discovery

	opts builderOptions
}

func NewBuilder(discovery registry.Discovery, opts ...BuilderOptions) resolver.Builder {
	b := &builder{
		discovery: discovery,
		opts: builderOptions{
			timeout:  10 * time.Second,
			insecure: false,
		},
	}

	for _, opt := range opts {
		opt(&b.opts)
	}
	return b
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var (
		err     error
		watcher registry.Watcher
	)

	done := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		watcher, err = b.discovery.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		if err != nil {
			return
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(b.opts.timeout):
		err = errors.New("discovery ccreate wather timeout")
	}

	if err != nil {
		cancel()
		return nil, err
	}
	r := &discoveryResolver{
		watcher:  watcher,
		cc:       cc,
		ctx:      ctx,
		cancel:   cancel,
		insecure: b.opts.insecure,
	}
	go r.watch()
	return r, nil
}

func (b *builder) Scheme() string {
	return name
}
