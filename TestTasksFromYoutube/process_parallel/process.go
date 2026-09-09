package process_parallel

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type processResult struct {
	res int
	err error
}

var errTimeout = errors.New("timeout")

func processData(ctx context.Context, v int) chan processResult {
	out := make(chan processResult)
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		close(ch)
	}()

	go func() {
		select {
		case <-ch:
			out <- processResult{res: v * 2, err: nil}
		case <-ctx.Done():
			out <- processResult{err: errTimeout}
		}
	}()

	return out
}

func Do(ctx context.Context) {
	in := make(chan int)
	out := make(chan int)

	go func() {
		defer close(in)

		//time.Sleep(10 * time.Second)
		for i := range 10 {
			select {
			case in <- i + 1:
			case <-ctx.Done():
				return
			}
		}
	}()

	start := time.Now()
	processParallel(ctx, in, out, 5)

	for v := range out {
		fmt.Println("v =", v)
	}

	fmt.Println("main duration:", time.Since(start))
}

func processParallel(ctx context.Context, in <-chan int, out chan<- int, numWorkers int) {
	wg := &sync.WaitGroup{}

	for range numWorkers {
		wg.Add(1)
		go worker(ctx, in, out, wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func worker(ctx context.Context, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case v, ok := <-in:
			if !ok {
				return
			}

			select {
			case r := <-processData(ctx, v):
				if r.err != nil {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- r.res:
				}
				out <- r.res
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
