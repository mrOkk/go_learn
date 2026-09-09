package main

import (
	"TestTasksFromYoutube/process_parallel"
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	process_parallel.Do(ctx)
}