package core

import (
	"context"
	"sync"
	"time"
)

func runTaskPeriodically(ctx context.Context, wg *sync.WaitGroup, task func(), interval time.Duration) {
	wg.Add(1)
	task()

	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			task()
		case <-ctx.Done():
			return
		}
	}
}
