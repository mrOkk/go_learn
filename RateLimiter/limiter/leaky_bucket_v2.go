package limiter

import (
	"context"
	"time"
)

type LeakyBucketV2 struct {
	tokenChan chan struct{}
}

func NewLeakyBucketV2(ctx context.Context, limit int, period time.Duration) *LeakyBucketV2 {
	lb := &LeakyBucketV2{
		tokenChan: make(chan struct{}, limit),
	}

	leakInterval := period.Nanoseconds() / int64(limit)
	go processLeak(ctx, lb, time.Duration(leakInterval))
	return lb
}

func processLeak(ctx context.Context, lb *LeakyBucketV2, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case <-lb.tokenChan:
			default:
			}
		}
	}
}

func (lb *LeakyBucketV2) Allow() bool {
	select {
	case lb.tokenChan <- struct{}{}:
		return true
	default:
		return false
	}
}
