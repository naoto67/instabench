package main

import (
	"context"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	cmd := NewRedisCommand()

	aggregator := NewAggregator()
	weighted := semaphore.NewWeighted(500)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	aggregator.Start()
	defer aggregator.Stop()
L:
	for {
		select {
		case <-ctx.Done():
			break L
		default:
			if weighted.TryAcquire(1) {
				go func() {
					defer weighted.Release(1)
					result, err := cmd.ExecContext(context.Background())
					if err != nil {
						aggregator.AddError(err)
					}
					aggregator.AddResult(result)
				}()
			}
		}
	}

	// fmt.Println(aggregator.td.Quantile(0.5))
	// fmt.Println(aggregator.td.Quantile(0.95))
}
