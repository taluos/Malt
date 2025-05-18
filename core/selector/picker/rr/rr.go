package rr

import (
	"context"
	"sync"

	"github.com/taluos/Malt/core/selector/picker/node/direct"

	"github.com/taluos/Malt/core/selector"
)

const (
	// Name is the name of the rr balancer.
	Name = "rr"
)

// rr is a round-robin balancer.

var _ selector.Balancer = &Balancer{} // Name is balancer name

// Balancer is a random balancer.
type Balancer struct {
	mu sync.Mutex
}

// New random a selector.
func New() selector.Selector {
	return NewBuilder().Build()
}

func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}

	var selected selector.WeightedNode

	p.mu.Lock()
	for _, node := range nodes {
		if selected == nil {
			selected = node
		}
	}

	d := selected.Pick()
	p.mu.Unlock()

	return selected, d, nil

}

// Builder is wrr builder
type Builder struct{}

// NewBuilder returns a selector builder with wrr balancer
func NewBuilder() selector.Builder {
	return &selector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Build creates Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{}
}
