package engine

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mohyehia/goku/internal/config"
	"github.com/mohyehia/goku/internal/metrics"
)

type HttpJob struct {
	JobID      string
	TargetURL  string
	HttpMethod string
}

type HttpResult struct {
	JobID      string
	StatusCode int
	Latency    time.Duration // latency will be in millisecond
	ErrorMsg   string
}

func Orchestrate(cfg *config.Config, aggregator *metrics.Aggregator) {
	jobs := make(chan HttpJob, 100)

	results := make(chan HttpResult, 100)

	// determine the numOfWorkers to be the min between 100 & cfg.concurrency to not exhaust the server
	numOfWorkers := min(cfg.Concurrency, 100)

	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var ctx context.Context
	var cancel context.CancelFunc

	if cfg.Requests > 0 {
		ctx, cancel = context.WithCancel(signalCtx)
	} else {
		ctx, cancel = context.WithTimeout(signalCtx, cfg.Duration)
	}
	defer cancel()

	workerWg := sync.WaitGroup{}
	consumerWg := sync.WaitGroup{}

	// spin up the workers based on concurrency
	for range numOfWorkers {
		workerWg.Go(func() {
			initializeWorker(ctx, httpClient, jobs, results)
		})
	}

	// spin up the consumers
	consumerWg.Go(func() {
		consumeJobs(results, aggregator)
	})

	// produce into jobs using a goroutine
	go produceJobs(ctx, cfg, jobs)

	go func() {
		workerWg.Wait()
		close(results)
	}()
	consumerWg.Wait()
}
