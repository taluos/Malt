package hash

import (
	"context"
	"hash/crc32"

	"github.com/taluos/Malt/core/selector/picker/node/direct"

	"github.com/taluos/Malt/core/selector"
)

const (
	// DefaultReplicas is the default number of replicas.
	Name = "hash"
)

type Balancer struct{}

func NewSelector() selector.Selector {
	return NewBuilder().Build()
}

func (p *Balancer) Pick(ctx context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	// 使用CRC32哈希算法计算方法名的哈希值，并取模得到节点索引
	methodName, ok := ctx.Value("FullMethod").(string)
	if !ok {
		return nil, nil, selector.ErrNoAvailable
	}
	hashVal := crc32.ChecksumIEEE([]byte(methodName))
	key := int(hashVal) % len(nodes)

	if key < 0 || key >= len(nodes) {
		key = 0
	}

	selected := nodes[key]
	d := selected.Pick()
	return selected, d, nil
}

func NewBuilder() selector.Builder {
	return &selector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder is wrr builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{}
}
