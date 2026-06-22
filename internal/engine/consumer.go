package engine

import (
	"github.com/mohyehia/goku/internal/metrics"
)

func consumeJobs(results <-chan HttpResult, aggregator *metrics.Aggregator) {
	for httpResult := range results {
		aggregator.Add(httpResult.Latency, httpResult.StatusCode, httpResult.ErrorMsg)
	}
}
