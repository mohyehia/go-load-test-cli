package engine

import (
	"context"
	"net/http"
)

func initializeWorker(ctx context.Context, httpClient *http.Client, headers map[string]string, jobs <-chan HttpJob, results chan<- HttpResult) {
	for job := range jobs {
		// call the API and return back the response
		httpResult := CallAPI(ctx, headers, httpClient, job)
		select {
		case <-ctx.Done():
			return
		case results <- httpResult:
		}
	}
}
