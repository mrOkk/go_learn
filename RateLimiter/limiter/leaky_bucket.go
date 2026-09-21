package limiter

import (
	"context"
	"sync/atomic"
	"time"
)

type LeakyBucket struct {
	limit uint32
	count atomic.Uint32
}

func NewLeakyBucket(ctx context.Context, limit int, interval time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		limit: uint32(limit),
	}
	go periodicLeak(ctx, lb, interval)
	return lb
}

func periodicLeak(ctx context.Context, lb *LeakyBucket, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count := lb.count.Load()
			if count > 0 {
				for !lb.count.CompareAndSwap(count, count-1) {
					count = lb.count.Load()
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (lb *LeakyBucket) Allow() bool {
	count := lb.count.Load()
	if count >= lb.limit {
		return false
	}

	for !lb.count.CompareAndSwap(count, count+1) {
		count = lb.count.Load()
	}
	return count < lb.limit
}
