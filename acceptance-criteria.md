# Load Tester CLI

## 📋 Acceptance Criteria (AC)

To consider this load-tester MVP (Minimum Viable Product) complete, it must satisfy the following criteria.

This phase should focus on a **clean, correct, and idiomatic Go CLI** that exercises:

- flag parsing and validation
- worker-pool concurrency
- request-scoped timeouts
- cancellation with `context.Context`
- safe result aggregation
- final summary reporting

## ✅ In Scope for This Phase (MVP)

### 1. Configuration & Input

- **AC 1.1:** The system **MUST** accept a target configuration including a valid URL and an HTTP Method (`GET`, `POST`, `PUT`, `DELETE`).
- **AC 1.2:** The system **MUST** allow the user to specify the load profile using *either* a fixed total number of requests **OR** a maximum execution duration.
- **AC 1.3:** The system **MUST** allow the user to configure the concurrency level (number of concurrent worker goroutines).
- **AC 1.4:** The system **MUST** reject invalid input before execution begins, including:
  - missing URL
  - malformed URL
  - unsupported HTTP method
  - invalid duration or timeout values
  - non-positive concurrency
  - non-positive request count
- **AC 1.5:** The system **MUST** enforce that exactly one execution mode is provided:
  - either `--requests`
  - or `--duration`
  - but not both
- **AC 1.6:** On invalid input, the CLI **MUST** print a clear validation error and exit with a non-zero status code.

### 2. Execution Engine

- **AC 2.1:** The system **MUST** execute requests concurrently using a bounded worker pool to prevent local resource exhaustion.
- **AC 2.2:** The system **MUST** support both execution strategies through a producer mechanism:
  - count-based execution for `--requests`
  - duration-based execution for `--duration`
- **AC 2.3:** The system **MUST** gracefully handle unexpected OS interrupt signals (`Ctrl+C`) by:
  - stopping the producer immediately
  - preventing new jobs from being enqueued
  - allowing active workers a clean window to finish in-flight requests
- **AC 2.4:** The system **MUST** enforce a timeout on each individual HTTP request so a hanging server does not freeze the worker pool.
- **AC 2.5:** The system **MUST** propagate cancellation through `context.Context` so producer, workers, and reporter all terminate cleanly.
- **AC 2.6:** The system **MUST NOT** leak goroutines after test completion or cancellation.

### 3. Metrics & Reporting

- **AC 3.1:** The system **MUST** capture execution metrics for every request, including:
  - HTTP status code (if a response was received)
  - latency in milliseconds
  - transport/connection error (if any)
- **AC 3.2:** The system **MUST** classify each result as either:
  - **Successful**: response received with status code `200-399`
  - **Failed**: status code `400-599` or any request/connection error
- **AC 3.3:** The system **MUST** stream real-time progress logs to the console without blocking worker execution.
- **AC 3.4:** Real-time progress reporting **MUST** be handled asynchronously by a dedicated reporting path or goroutine.
- **AC 3.5:** Upon completion or cancellation, the system **MUST** print a final summary report containing at minimum:
  - total requests attempted (requests actually started by workers)
  - successful count
  - failed count
  - average latency
  - total elapsed time
  - requests per second (RPS / throughput)
- **AC 3.6:** The final summary **SHOULD** also include:
  - minimum latency
  - maximum latency
  - status code distribution
  - total connection/timeout errors
- **AC 3.7:** If the run is interrupted, the system **MUST** still print a partial final summary for all completed requests.

### 4. CLI UX & Output

- **AC 4.1:** The CLI **MUST** expose the documented flags exactly as specified for the MVP.
- **AC 4.2:** The CLI **MUST** provide sensible defaults for:
  - HTTP method = `GET`
  - concurrency = `10`
  - request timeout = `5s`
- **AC 4.3:** The CLI **MUST** print human-readable output to standard output/error suitable for local terminal use.
- **AC 4.4:** Validation and runtime error messages **SHOULD** be concise and actionable.

### 5. Code Quality & Learning Goals

- **AC 5.1:** The implementation **MUST** separate concerns clearly between:
  - CLI/config parsing
  - job production
  - worker execution
  - result aggregation
  - reporting
- **AC 5.2:** Shared metrics/state **MUST** be aggregated safely under concurrency.
- **AC 5.3:** The codebase **SHOULD** be structured in a way that allows unit testing of:
  - input validation
  - execution mode selection
  - result aggregation
  - summary calculation

## ⚙️ User-Allowed Options (The Command Flags)

When the user runs `goku` from their terminal, they should be able to pass these exact flags.

| **Flag / Option**     | **Type** | **Description**                                                                       | **Mapping to Code**                                                         |
|-----------------------|----------|---------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| `-u`, `--url`         | `string` | The target API endpoint to test (e.g., `https://api.example.com/v1/users`).           | **Required.** Passed to the HTTP client inside the worker request flow.     |
| `-m`, `--method`      | `string` | The HTTP verb to use. Defaults to `GET`.                                              | Passed to `http.NewRequestWithContext`.                                     |
| `-c`, `--concurrency` | `int`    | The number of concurrent virtual users running simultaneously. Defaults to `10`.      | Maps directly to the worker pool size.                                      |
| `-n`, `--requests`    | `int`    | Total number of requests to fire. *Required if duration is not provided.*             | Maps to the count-based producer.                                           |
| `-d`, `--duration`    | `string` | How long to run the test (e.g., `30s`, `2m`). *Required if requests is not provided.* | Parsed via `time.ParseDuration` and converted into a `context.WithTimeout`. |
| `-t`, `--timeout`     | `string` | Hard timeout limit for any single HTTP request. Defaults to `5s`.                     | Wrapped into each worker's individual request context.                      |

## 🔄 Choosing the Execution Strategy

To handle the user choice between a **count-based test** (`-n`) and a **duration-based test** (`-d`), the orchestration layer should evaluate which flag was provided and assign the corresponding producer strategy.

- **If `-n 50000` is passed:** `goku` runs a count-based sequential producer loop until `50000` jobs have been submitted or cancellation occurs.
- **If `-d 60s` is passed:** `goku` attaches a 60-second timeout to the run context and uses a duration-based producer loop that stops when the context is done.

## 🧭 Recommended Actions for This Phase

Build the MVP in this order so the project stays small, testable, and aligned with your Go learning goals:

1. **Create a config layer**
   - Parse flags into a config struct.
   - Validate mutual exclusivity of `--requests` and `--duration`.
   - Normalize and validate method, URL, concurrency, and timeout values.

2. **Implement the execution orchestration**
   - Create the root context.
   - Attach timeout when duration mode is used.
   - Listen for OS interrupts and trigger cancellation.

3. **Implement a bounded worker pool**
   - Use a jobs channel for request work.
   - Spin up `N` workers from the configured concurrency value.
   - Ensure workers stop cleanly when jobs close or context is canceled.

4. **Implement producer strategies**
   - Count-based producer for `--requests`.
   - Duration-based producer for `--duration`.
   - Stop producing immediately on cancellation.

5. **Capture and aggregate results safely**
   - Send worker results through a results channel.
   - Aggregate counts and latency statistics in one place.
   - Avoid workers mutating shared state directly.

6. **Add asynchronous progress reporting**
   - Use a separate goroutine or reporting path.
   - Print periodic progress without blocking workers.

7. **Print a final summary**
   - Show totals, successes, failures, average latency, elapsed time, and RPS.
   - Still print a partial summary when the run is canceled.

8. **Add focused tests**
   - Validation tests for flags/config.
   - Aggregation tests for summary calculations.
   - Small orchestration tests for count-based vs duration-based mode.

## 🚀 Explicitly Out of Scope for This MVP

These are valuable next steps, but they should not block completion of this phase:

- custom request headers
- request body support for `POST` / `PUT`
- JSON or CSV output modes
- latency percentiles (`p50`, `p90`, `p95`, `p99`)
- rate limiting (requests per second caps)
- Prometheus metrics endpoint
- OpenTelemetry tracing
- advanced HTTP transport tuning
