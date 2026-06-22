package engine

import (
	"context"
	"io"
	"net/http"
	"time"
)

func CallAPI(ctx context.Context, httpClient *http.Client, job HttpJob) HttpResult {
	req, err := http.NewRequestWithContext(ctx, job.HttpMethod, job.TargetURL, nil)
	if err != nil {
		return HttpResult{
			JobID:    job.JobID,
			ErrorMsg: err.Error(),
		}
	}
	startTime := time.Now()
	res, err := httpClient.Do(req)
	latency := time.Since(startTime)
	// round the latency to get only the first 2 digits after the decimal point. EX: 440.44
	latency = latency.Round(10 * time.Microsecond)
	if err != nil {
		return HttpResult{
			JobID:    job.JobID,
			Latency:  latency,
			ErrorMsg: err.Error(),
		}
	}
	defer func() {
		// the below function is to reuse the same network connection for other workers
		// This is just for enabling  HTTP Connection Re-use (Keep-Alive)
		_, _ = io.Copy(io.Discard, res.Body)
		err := res.Body.Close()
		if err != nil {
			return
		}
	}()

	return HttpResult{
		JobID:      job.JobID,
		StatusCode: res.StatusCode,
		Latency:    latency,
		ErrorMsg:   "",
	}
}
