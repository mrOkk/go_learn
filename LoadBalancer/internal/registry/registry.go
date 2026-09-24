package registry

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"
)

type Registry struct {
	mu          sync.RWMutex
	backendsMap map[string]int
	backends    []*Backend
	ttl         time.Duration
	count       int
}

func NewRegistry(ttl time.Duration) *Registry {
	return &Registry{
		backendsMap: make(map[string]int),
		backends:    make([]*Backend, 0, 4),
		ttl:         ttl,
	}
}

func (r *Registry) Register(address string) (string, time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	backend := &Backend{
		Id:       fmt.Sprintf("srv%d", r.count),
		Address:  address,
		LastSeen: time.Now(),
	}
	r.backendsMap[backend.Id] = len(r.backends)
	r.backends = append(r.backends, backend)
	r.count++
	return backend.Id, r.ttl
}

func (r *Registry) Deregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	bIndex, found := r.backendsMap[id]
	if !found {
		return
	}
	r.backends = slices.Delete(r.backends, bIndex, bIndex+1)
	delete(r.backendsMap, id)
	for i := bIndex; i < len(r.backends); i++ {
		b := r.backends[i]
		r.backendsMap[b.Id] = i
	}
}

func (r *Registry) Heartbeat(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	bIndex, ok := r.backendsMap[id]
	if !ok {
		return fmt.Errorf("no backend for id %s", id)
	}
	backend := r.backends[bIndex]
	backend.LastSeen = time.Now()
	return nil
}

func (r *Registry) Snapshot() []*Backend {
	r.mu.RLock()
	defer r.mu.RUnlock()
	backends := make([]*Backend, len(r.backends))
	copy(backends, r.backends)
	return backends
}

func (r *Registry) SweepJob(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.sweep(time.Now())
		}
	}
}

func (r *Registry) sweep(now time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	anyDeleted := false
	for i, b := range slices.Backward(r.backends) {
		if now.Sub(b.LastSeen) > r.ttl {
			r.backends = slices.Delete(r.backends, i, i+1)
			delete(r.backendsMap, b.Id)
			anyDeleted = true
		}
	}
	if !anyDeleted {
		return
	}
	for i, b := range r.backends {
		r.backendsMap[b.Id] = i
	}
}
