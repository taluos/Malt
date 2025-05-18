package chash

import (
	"context"
	"fmt"
	"hash/crc32"
	"sort"
	"sync"

	"Malt/core/selector"
	"Malt/core/selector/picker/node/direct"
)

const (
	// Name is the name of the random balancer.
	Name = "chash"
)

type Balancer struct {
	ring *hashRing
}

type hashRing struct {
	replicas int
	keys     []uint32
	hashMap  map[uint32]selector.WeightedNode
	mu       sync.RWMutex
}

func NewSelector() selector.Selector {
	return NewBuilder().Build()
}

func (p *Balancer) Pick(ctx context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}

	p.ring = &hashRing{
		replicas: 5,
		keys:     make([]uint32, 0, len(nodes)),
		hashMap:  make(map[uint32]selector.WeightedNode),
	}

	p.ring.mu.Lock()
	for _, node := range nodes {
		for i := 0; i < p.ring.replicas; i++ {
			key := fmt.Sprintf("%s#%d", node.ID(), i)
			h := crc32.ChecksumIEEE([]byte(key))
			p.ring.keys = append(p.ring.keys, h)
			p.ring.hashMap[h] = node
		}
	}
	sort.Slice(p.ring.keys, func(i, j int) bool { return p.ring.keys[i] < p.ring.keys[j] })
	p.ring.mu.Unlock()

	key, ok := ctx.Value("FullMethod").(string)
	if !ok {
		return nil, nil, selector.ErrNoAvailable
	}

	h := crc32.ChecksumIEEE([]byte(key))

	p.ring.mu.RLock()
	idx := sort.Search(len(p.ring.keys), func(i int) bool {
		return p.ring.keys[i] >= h
	})
	if idx == len(p.ring.keys) {
		idx = 0
	}
	selected := p.ring.hashMap[p.ring.keys[idx]]
	p.ring.mu.RUnlock()

	if selected == nil {
		return nil, nil, selector.ErrNoAvailable
	}

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
