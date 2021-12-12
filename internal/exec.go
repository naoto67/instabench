package internal

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/naoto67/instabench"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/semaphore"
)

func exec(cliContext *cli.Context, executor instabench.Executor) error {
	weighted := semaphore.NewWeighted(thread)
	wg := &sync.WaitGroup{}

	var (
		err error
		ctx = cliContext.Context
	)
	ctx, err = executor.Setup(ctx)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(cliContext.Context, time.Duration(duration)*time.Second)
	defer cancel()

	aggregator := instabench.NewAggregator()
	aggregator.Start()
	defer aggregator.Stop()
L:
	for {
		select {
		case <-ctx.Done():
			break L
		default:
			if weighted.TryAcquire(1) {
				wg.Add(1)
				go func() {
					defer weighted.Release(1)
					defer wg.Done()
					var (
						gctx = context.Background()
						gErr error
					)
					if preparer, ok := executor.(instabench.Preparer); ok {
						gctx, gErr = preparer.PrepareContext(gctx)
						if gErr != nil {
							return
						}
					}
					began := time.Now()
					gctx, gErr = executor.ExecContext(gctx)
					result := &instabench.Result{
						Timestamp: began,
						Latency:   time.Since(began),
						Err:       gErr,
					}
					aggregator.AddResult(result)
				}()
			}
		}
	}

	wg.Wait()
	report := &instabench.Report{
		Config: map[string]interface{}{
			"service":       service,
			"duration(sec)": duration,
			"threads":       thread,
		},
		Results: aggregator.ResultsSlice,
	}
	var writer io.Writer = os.Stdout
	if path := cliContext.String(outputJsonPath); path != "" {
		var file *os.File
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
	}
	if v, ok := executor.(instabench.ExtraConfigExporter); ok {
		report.ExtraConfig = v.ExtraConfig()
	}
	err = json.NewEncoder(writer).Encode(report)
	if err != nil {
		return err
	}
	return nil
}
