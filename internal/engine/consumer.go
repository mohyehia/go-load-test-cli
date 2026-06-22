package engine

import (
	"fmt"

	"github.com/mohyehia/goku/internal/metrics"
)

func consumeJobs(results <-chan HttpResult, aggregator *metrics.Aggregator) {
	for httpResult := range results {
		fmt.Printf("Adding httpResult into the metrics aggregator: %v\n", httpResult)
		aggregator.Add(httpResult.Latency, httpResult.StatusCode, httpResult.ErrorMsg)
	}
}
