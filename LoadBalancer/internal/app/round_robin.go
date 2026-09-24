package app

import (
	"LoadBalancer/internal/registry"
	"fmt"
	"sync/atomic"
)

type RoundRobin struct {
	reg     BackendContainer
	counter atomic.Uint64
}

func NewRoundRobin(reg BackendContainer) *RoundRobin {
	return &RoundRobin{reg: reg}
}

func (r *RoundRobin) Pick() (*registry.Backend, error) {
	backs := r.reg.Snapshot()
	backCount := len(backs)

	if backCount == 0 {
		return nil, fmt.Errorf("no available backend")
	}

	idx := r.counter.Add(1) % uint64(backCount)
	return backs[idx], nil
}
