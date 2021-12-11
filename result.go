package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/influxdata/tdigest"
)

type Aggregator struct {
	mux *sync.Mutex

	td         *tdigest.TDigest
	totalCount int64
	count      int64
	errs       []error

	closeCh chan struct{}
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		mux:  &sync.Mutex{},
		td:   tdigest.New(),
		errs: []error{},

		closeCh: make(chan struct{}),
	}
}

func (a *Aggregator) AddResult(result *Result) {
	if result == nil {
		return
	}
	a.mux.Lock()
	defer a.mux.Unlock()
	dur := float64(result.Latency) / float64(time.Millisecond)
	a.td.Add(dur, 1)
	a.count++
	a.totalCount++
}

func (a *Aggregator) AddError(err error) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.errs = append(a.errs, err)
}

func (a *Aggregator) Start() {
	go a.run()
}

func (a *Aggregator) Stop() {
	close(a.closeCh)

	fmt.Println(len(a.errs))
}

func (a *Aggregator) run() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-a.closeCh:
			return
		case <-ticker.C:
			a.mux.Lock()
			cnt := a.count
			// totalCnt := a.totalCount
			a.count = 0
			a.mux.Unlock()
			fmt.Printf("%s RPS: %d\n", time.Now().Format(timeFormat), cnt)
		}
	}
}

type Result struct {
	Timestamp time.Time
	Latency   time.Duration
}
