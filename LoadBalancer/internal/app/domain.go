package app

import "LoadBalancer/internal/registry"

type Picker interface {
	Pick() (*registry.Backend, error)
}

type BackendContainer interface {
	Snapshot() []*registry.Backend
}
