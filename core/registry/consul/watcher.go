package consul

import (
	"context"

	"Malt/core/registry"
)

// Use assertions to ensure that the watcher struct implements the registry.Watcher interface.
var _ registry.Watcher = (*watcher)(nil)

type watcher struct {
	event chan struct{}
	set   *serviceSet

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *watcher) Next() (services []*registry.ServiceInstance, err error) {
	if err = w.ctx.Err(); err != nil {
		return
	}

	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
		return
	case <-w.event:
		// this means we accept a message from brocast ( service.go/broadcast() )
		// or form watcher ( registry.go/watch() )
		//
		// if w.event is not nil, it means that there is a new service instance
		// or a serveive instance is updated or deleted
		// if not done, coutinue to get service instances
	}

	ss, ok := w.set.services.Load().([]*registry.ServiceInstance)
	if ok {
		services = append(services, ss...)
	}
	return
}

func (w *watcher) Stop() error {
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
		w.set.delete(w)
	}
	return nil
}
