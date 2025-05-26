// In this file we define discorveryResolver struct and its methods,
// to replace the grpc resolver.Resolver interface by your own.
// In specific, we add a watch() method to watch the discovery endpoint,
// and a update() methd to update the state of resolver.Resolver interface by the watcher.
//
// A discoveryResolver instance is created at Build() method in builder.go.
package discovery

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/taluos/Malt/core/registry"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/log"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

var _ resolver.Resolver = (*discoveryResolver)(nil)

type discoveryResolver struct {
	watcher registry.Watcher
	cc      resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	insecure bool
}

// update: main logic of resolver
// 以 watcher （观察者）模式实现Resolver接口
func (r *discoveryResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		ins, err := r.watcher.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Errorf("[resolver] Failed to watch discorvery endpoint: %v", err)
			time.Sleep(watchSleep)
			continue
		}
		r.update(ins)
	}
}

func (r *discoveryResolver) update(ins []*registry.ServiceInstance) {
	address := make([]resolver.Address, 0)
	endpointes := make(map[string]struct{})
	for _, in := range ins {
		endpoint, err := ParseEndpoint(in.Endpoints, "grpc", !r.insecure)
		if err != nil {
			log.Errorf("[resolver] Failed to parse disconvery endpoint: %v", err)
			continue
		}
		if endpoint == "" {
			continue
		}

		// filter redundant endpoints
		if _, ok := endpointes[endpoint]; ok {
			continue
		}
		endpointes[endpoint] = struct{}{}

		addr := resolver.Address{
			ServerName: in.Name,
			Addr:       endpoint,
			Attributes: parseAttributes(in.Metadata),
		}
		addr.Attributes = addr.Attributes.WithValue("rawServiceInstance", in)
		address = append(address, addr)
	}
	if len(address) == 0 {
		log.Warnf("[resolver] No available endpoint")
		return
	}
	err := r.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	if err != nil {
		log.Errorf("[resolver] Failed to update state: %v", err)
	}
	b, _ := json.Marshal(ins)
	log.Infof("[resolver] Update state: %s", string(b))
}

// 实现 resolver.Resolver 接口
func (r *discoveryResolver) ResolveNow(rn resolver.ResolveNowOptions) {}

func (r *discoveryResolver) Close() {
	r.cancel()
	err := r.watcher.Stop()
	if err != nil {
		log.Errorf("[resolver] Failed to stop discovery watcher: %v", err)
	}
}

func parseAttributes(md map[string]string) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range md {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}

func NewEndpoint(schema, host string, isSecure bool) *url.URL {
	var query string
	if isSecure {
		query = "isSecure=true"
	}
	return &url.URL{
		Scheme:   schema,
		Host:     host,
		RawQuery: query,
	}
}

func ParseEndpoint(endpoints []string, schema string, isSecure bool) (string, error) {
	for _, eendpoint := range endpoints {
		u, err := url.Parse(eendpoint)
		if err != nil {
			return "", err
		}
		if u.Scheme == schema {
			if IsSecure(u) == isSecure { // || isSecure
				return u.Host, nil
			}
		}
	}
	return "", nil
}

func IsSecure(u *url.URL) bool {
	ok, err := strconv.ParseBool(u.Query().Get("isSecure"))
	if err != nil {
		return false
	}
	return ok
}
