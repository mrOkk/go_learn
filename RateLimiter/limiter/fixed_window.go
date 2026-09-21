package limiter

import (
	"context"
	"sync/atomic"
	"time"
)

type FixedWindow struct {
	limit uint32
	count atomic.Uint32
}

func NewFixedWindow(ctx context.Context, limit uint32) *FixedWindow {
	fw := &FixedWindow{
		limit: limit,
	}

	go resetWindowJob(ctx, fw)
	return fw
}

func resetWindowJob(ctx context.Context, w *FixedWindow) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.count.Store(0)
		}
	}
}

func (fw *FixedWindow) Allow() bool {
	count := fw.count.Load()
	if count >= fw.limit {
		return false
	}

	for !fw.count.CompareAndSwap(count, count+1) {
		count = fw.count.Load()
	}
	return count < fw.limit
}
