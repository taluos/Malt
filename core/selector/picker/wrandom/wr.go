package wrandom

import (
	"context"
	"math/rand"
	"sync"

	"github.com/taluos/Malt/core/selector/picker/node/direct"

	"github.com/taluos/Malt/core/selector"
)

const (
	// Name is the name of wrr balancer.
	Name = "wrandom"
)

var _ selector.Balancer = &Balancer{} // Name is balancer name

type Balancer struct {
	mu            sync.RWMutex
	currentWeight map[string]wrange
}

type wrange struct {
	Start float64
	End   float64
}

func NewSelector() selector.Selector {
	return NewBuilder().Build()
}

func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}

	var totalWeight float64
	var selected selector.WeightedNode

	p.mu.Lock()
	for _, node := range nodes {

		weight := node.Weight()
		cwt := p.currentWeight[node.Address()]

		cwt.Start = totalWeight
		totalWeight += weight
		cwt.End = totalWeight

		p.currentWeight[node.Address()] = cwt
	}
	p.mu.Unlock()

	cur := rand.Float64() * totalWeight

	p.mu.RLock()

	for _, node := range nodes {
		cwt := p.currentWeight[node.Address()]
		if cur >= cwt.Start && cur < cwt.End {
			selected = node
			break
		}
	}

	p.mu.RUnlock()

	if selected == nil {
		return nil, nil, selector.ErrNoAvailable
	}

	d := selected.Pick()

	return selected, d, nil

}

// NewBuilder returns a selector builder with wrr balancer
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
	return &Balancer{currentWeight: make(map[string]wrange)}
}
