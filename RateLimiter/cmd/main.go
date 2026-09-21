package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"RateLimiter/limiter"
)

func main() {
	const maxRequestsCount = 60
	const minRequestCount = 30
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := limiter.NewFixedWindow(ctx, 100)

	for i := range 20 {
		approved := 0
		rc := minRequestCount + rand.N(maxRequestsCount-minRequestCount)
		for range rc {
			if l.Allow() {
				approved++
			}
		}
		fmt.Println("iteration:", i, "| approved count:", approved, "/", rc)
		time.Sleep(350 * time.Millisecond)
	}
}
