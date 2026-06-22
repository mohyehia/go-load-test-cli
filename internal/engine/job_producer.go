package engine

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mohyehia/goku/internal/config"
)

func produceJobs(ctx context.Context, cfg *config.Config, jobs chan<- HttpJob) {
	if cfg.Requests > 0 {
		// profile => Count-based
		produceJobsBasedOnCount(ctx, cfg.Requests, cfg.URL, cfg.Method, jobs)
	} else {
		// profile => Duration-based
		produceJobsBasedOnDuration(ctx, cfg.URL, cfg.Method, jobs)
	}
}

func produceJobsBasedOnCount(ctx context.Context, totalRequests int, URL string, httpMethod string, jobs chan<- HttpJob) {
	defer close(jobs)
	for range totalRequests {
		job := HttpJob{
			JobID:      uuid.New().String(),
			TargetURL:  URL,
			HttpMethod: httpMethod,
		}
		select {
		case <-ctx.Done():
			fmt.Printf("🛑 Job Producer shutting down gracefully, exiting...\n")
			return
		case jobs <- job:
		}
	}
}

func produceJobsBasedOnDuration(ctx context.Context, URL string, httpMethod string, jobs chan<- HttpJob) {
	defer close(jobs)
	for {
		job := HttpJob{
			JobID:      uuid.New().String(),
			TargetURL:  URL,
			HttpMethod: httpMethod,
		}
		select {
		case <-ctx.Done():
			fmt.Printf("🛑 Job Producer shutting down gracefully, exiting...\n")
			return
		case jobs <- job:
		}
	}
}
