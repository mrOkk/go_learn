package app

import (
	"context"
	"time"
)

type Registry interface {
	Register(address string) (string, time.Duration)
	Heartbeat(id string) error
	Deregister(id string)
}

type BackendCaller interface {
	Shorten(ctx context.Context, address, url string) (string, error)
	Resolve(ctx context.Context, address, code string) (string, error)
}
