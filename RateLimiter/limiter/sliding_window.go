package limiter

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	limit    int
	interval time.Duration
	mu       sync.Mutex

	windowStart time.Time
	prevCount   int
	count       int
}

func NewSlidingWindow(limit int, interval time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:       limit,
		interval:    interval,
		windowStart: time.Now(),
	}
}

func (s *SlidingWindow) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	newPeriod := s.windowStart.Add(s.interval)
	if now.After(newPeriod) {
		s.windowStart = newPeriod
		s.prevCount = s.count
		s.count = 0
	}

	interval := float64(s.interval)
	currentCount := float64(s.count)
	prevCount := float64(s.prevCount)
	elapsed := now.Sub(s.windowStart).Seconds()
	count := (prevCount * (interval - elapsed) / interval) + currentCount

	if int(count) >= s.limit {
		return false
	}

	s.count++
	return true
}
