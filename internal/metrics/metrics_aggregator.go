package metrics

import (
	"fmt"
	"sync"
	"time"
)

type Aggregator struct {
	metrics *Metrics
	sync.Mutex
}

type Metrics struct {
	totalCount        int
	successCount      int // no full storage needed
	failedCount       int // no full storage needed
	latencySum        time.Duration
	averageLatency    time.Duration // no full storage needed
	totalTime         time.Duration // no full storage needed
	requestsPerSecond float64       // no full storage needed
	minLatency        time.Duration // no full storage needed
	maxLatency        time.Duration // no full storage needed
	totalErrorCount   int
	startTime         time.Time
	endTime           time.Time
}

func NewAggregator() *Aggregator {
	m := &Metrics{
		startTime: time.Now(),
	}
	return &Aggregator{
		metrics: m,
	}
}

func (a *Aggregator) Add(latency time.Duration, statusCode int, errorMsg string) {
	a.Lock()
	defer a.Unlock()
	a.metrics.totalCount++
	if statusCode >= 200 && statusCode < 400 {
		a.metrics.successCount++
	} else {
		a.metrics.failedCount++
	}
	if errorMsg != "" {
		a.metrics.totalErrorCount++
		return // Do not calculate latency for failed/aborted network calls
	}

	a.metrics.latencySum += latency

	if a.metrics.totalCount == 1 {
		// First request
		a.metrics.minLatency = latency
		a.metrics.maxLatency = latency
	} else {
		a.metrics.minLatency = min(latency, a.metrics.minLatency)
		a.metrics.maxLatency = max(latency, a.metrics.maxLatency)
	}
}

func (a *Aggregator) Aggregate() {

	a.Lock()
	defer a.Unlock()

	if a.metrics.totalCount == 0 {
		fmt.Println("\n================ GOKU RESULTS ================")
		fmt.Println("No requests were successfully processed.")
		fmt.Println("==============================================")
		return
	}

	a.metrics.endTime = time.Now()
	a.metrics.totalTime = time.Since(a.metrics.startTime)
	a.metrics.requestsPerSecond = float64(a.metrics.totalCount) / a.metrics.totalTime.Seconds()

	// Calculate Average Latency based ONLY on actual network round-trips completed
	validRequests := a.metrics.totalCount - a.metrics.totalErrorCount
	if validRequests > 0 {
		a.metrics.averageLatency = a.metrics.latencySum / time.Duration(validRequests)
	}

	fmt.Println("\n================ GOKU RESULTS ================")
	fmt.Printf("Total execution time:   %v\n", a.metrics.totalTime.Round(time.Millisecond))
	fmt.Printf("Throughput (RPS):       %.2f req/sec\n", a.metrics.requestsPerSecond)
	fmt.Printf("Total requests count:   %d\n", a.metrics.totalCount)
	fmt.Printf("✅ Successful (2xx - 3xx):    %d\n", a.metrics.successCount)
	fmt.Printf("❌ Failed (Non-2xx):    %d\n", a.metrics.failedCount-a.metrics.totalErrorCount) // Actual HTTP status code failures
	fmt.Printf("⚠️ Network/OS Errors:   %d\n", a.metrics.totalErrorCount)
	fmt.Println("----------------------------------------------")
	fmt.Printf("Latency Sum:            %v\n", a.metrics.latencySum.Round(time.Millisecond))
	fmt.Printf("Avg Request Latency:    %v\n", a.metrics.averageLatency.Round(time.Microsecond))
	fmt.Printf("Min Request Latency:    %v\n", a.metrics.minLatency.Round(time.Microsecond))
	fmt.Printf("Max Request Latency:    %v\n", a.metrics.maxLatency.Round(time.Microsecond))
	fmt.Println("==============================================")
}

func (a *Aggregator) PrintProgress() {
	a.Lock()
	defer a.Unlock()
	// Guard against division by zero if it ticks before the first request finishes
	// \r resets the cursor to the start of the line.
	if a.metrics.totalCount == 0 {
		fmt.Print("\r⏳ [Goku] Preparing workers...")
		return
	}
	totalTime := time.Since(a.metrics.startTime)

	// calculate requests per second
	rps := float64(a.metrics.totalCount) / totalTime.Seconds()
	// The trailing space clears out any leftover characters from previous longer lines.
	fmt.Printf("\r 🚀 [Goku] Running... Requests: %d | Success: %d | Failed: %d | Current RPS: %.2f	",
		a.metrics.totalCount,
		a.metrics.successCount,
		a.metrics.failedCount,
		rps,
	)
}
