package consul

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/taluos/Malt/core/registry"
)

type serviceSet struct {
	registry    *Registry
	serviceName string
	watcher     map[*watcher]struct{}
	ref         atomic.Int32
	services    *atomic.Value
	lock        sync.RWMutex

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

// broadcast sends a message to all watchers.
func (s *serviceSet) broadcast(ss []*registry.ServiceInstance) {
	// actomic store the service instances
	s.services.Store(ss)
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}

// delete deletes the watcher.
func (s *serviceSet) delete(w *watcher) {
	s.lock.Lock()
	delete(s.watcher, w)
	s.lock.Unlock()
	s.registry.tryDelete(s)
}
