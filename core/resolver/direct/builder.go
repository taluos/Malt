package direct

import (
	"Malt/pkg/errors"

	"strings"

	"google.golang.org/grpc/resolver"
)

var _ resolver.Builder = (*directBuilder)(nil)

// directBuilder 是直接服务发现的 resolver builder 实现
type directBuilder struct{}

// NewDirectBuilder create a new directBuilder which is used to build direct resolvers.
// example:
//
//	direct://<authority>/127.0.0.1:8080
func NewDirectBuilder() *directBuilder {
	return &directBuilder{}
}

// Build create a new direct resolver
func (b *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 解析 target 地址
	endpoints := strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",")
	if len(endpoints) == 0 || endpoints[0] == "" {
		return nil, errors.New("direct resolver: no endpoints provided")
	}

	// 创建地址列表
	addrs := make([]resolver.Address, 0)
	for _, endpoint := range endpoints {
		if endpoint == "" {
			continue
		}
		addrs = append(addrs, resolver.Address{
			Addr: endpoint,
		})
	}

	// grpc 连接
	err := cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
	if err != nil {
		return nil, errors.Wrap(err, "direct resolver: update state failed")
	}

	// 返回一个空的 resolver，因为直接模式不需要监视服务变更
	return &directResolver{cc: cc}, nil
}

// Scheme 返回 direct resolver 的 scheme
func (b *directBuilder) Scheme() string {
	return name
}
