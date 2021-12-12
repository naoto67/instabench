package instabench

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/influxdata/tdigest"
)

type Aggregator struct {
	mux          *sync.Mutex
	ResultsSlice []*Results

	closeCh chan struct{}
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		mux:          &sync.Mutex{},
		ResultsSlice: make([]*Results, 0, 100), // TODO 100: Duration

		closeCh: make(chan struct{}),
	}
}

func (a *Aggregator) AddResult(result *Result) {
	if result == nil {
		return
	}
	a.mux.Lock()
	defer a.mux.Unlock()
	idx := len(a.ResultsSlice) - 1
	a.ResultsSlice[idx].results = append(a.ResultsSlice[idx].results, result)
}

func (a *Aggregator) Start() {
	a.ResultsSlice = append(a.ResultsSlice, &Results{
		results:      make([]*Result, 0, 100),
		startTime:    time.Now(),
		sumLatencyMs: big.NewFloat(0),
		td:           tdigest.New(),
	})
	go a.run()
}

func (a *Aggregator) Stop() {
	close(a.closeCh)
}

func (a *Aggregator) run() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-a.closeCh:
			return
		case <-ticker.C:
			a.mux.Lock()
			idx := len(a.ResultsSlice) - 1
			a.ResultsSlice = append(a.ResultsSlice, &Results{
				results:      make([]*Result, 0, len(a.ResultsSlice[idx].results)),
				startTime:    time.Now(),
				sumLatencyMs: big.NewFloat(0),
				td:           tdigest.New(),
			})
			a.ResultsSlice[idx].Bundle()
			a.mux.Unlock()
			fmt.Println(a.ResultsSlice[idx])
		}
	}
}

type Result struct {
	Timestamp time.Time
	Latency   time.Duration

	Err error
}

type Results struct {
	results   []*Result
	startTime time.Time

	td           *tdigest.TDigest
	sumLatencyMs *big.Float
}

func (r *Results) Bundle() {
	td := tdigest.New()
	for i := range r.results {
		dur := float64(r.results[i].Latency) / float64(time.Millisecond)
		td.Add(dur, 1)
		r.sumLatencyMs.Add(r.sumLatencyMs, big.NewFloat(dur))
	}
	r.td = td
}

func (r *Results) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, 9)
	m["rps"] = len(r.results)
	m["start_time"] = r.startTime
	if len(r.results) > 0 {
		avg, _ := r.sumLatencyMs.Quo(r.sumLatencyMs, big.NewFloat(float64((len(r.results))))).Float64()
		if avg != 0.0 {
			m["avg(ms)"] = avg
			m["min(ms)"] = r.td.Quantile(0.0)
			m["p50(ms)"] = r.td.Quantile(0.5)
			m["p90(ms)"] = r.td.Quantile(0.9)
			m["p95(ms)"] = r.td.Quantile(0.95)
			m["p99(ms)"] = r.td.Quantile(0.99)
			m["max(ms)"] = r.td.Quantile(1.0)
		}
	}
	return json.Marshal(m)
}

func (r *Results) String() string {
	return fmt.Sprintf(`
StartTime:       %s,
RPS:             %d,
P50:             %f(ms),
P90:             %f(ms),
P95:             %f(ms),
P90:             %f(ms),
`,
		r.startTime.Format(timeFormat),
		len(r.results),
		r.td.Quantile(float64(0.5)),
		r.td.Quantile(float64(0.9)),
		r.td.Quantile(float64(0.95)),
		r.td.Quantile(float64(0.99)),
	)
}
