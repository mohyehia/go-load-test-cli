package config

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	URL         string
	Method      string
	Requests    int
	Concurrency int
	Duration    time.Duration
	Timeout     time.Duration
}

func ParseFlags() (*Config, error) {
	config := &Config{}
	// Bind short and long flags to the exact same struct pointers
	flag.StringVar(&config.URL, "u", "", "Target URL to load test (Required)")
	flag.StringVar(&config.URL, "url", "", "Target URL to load test (Required)")

	flag.StringVar(&config.Method, "m", "GET", "HTTP Method to use")
	flag.StringVar(&config.Method, "method", "GET", "HTTP Method to use")

	flag.IntVar(&config.Requests, "n", 0, "Number of requests")
	flag.IntVar(&config.Requests, "requests", 0, "Number of requests")

	flag.IntVar(&config.Concurrency, "c", 10, "Number of concurrent requests")
	flag.IntVar(&config.Concurrency, "concurrency", 10, "Number of concurrent requests")

	var duration string
	flag.StringVar(&duration, "d", "", "Duration of the test (e.g., 30s, 1m)")
	flag.StringVar(&duration, "duration", "", "Duration of the test (e.g., 30s, 1m)")

	var timeOut string
	flag.StringVar(&timeOut, "t", "5s", "Timeout limit per individual HTTP request")
	flag.StringVar(&timeOut, "timeout", "5s", "Timeout limit per individual HTTP request")

	// Parse parameters out of os.Args
	flag.Parse()

	// update method to be uppercase
	config.Method = strings.ToUpper(strings.TrimSpace(config.Method))

	if strings.TrimSpace(duration) != "" {
		d, err := time.ParseDuration(duration)
		if err != nil {
			return nil, fmt.Errorf("invalid duration format '%s': %v", duration, err)
		}
		config.Duration = d
	}

	if strings.TrimSpace(timeOut) != "" {
		d, err := time.ParseDuration(timeOut)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout format '%s': %v", timeOut, err)
		}
		config.Timeout = d
	}

	return config, nil
}

func (c *Config) Validate() error {
	if err := validateURL(c.URL); err != nil {
		return err
	}

	if err := validateHttpMethod(c.Method); err != nil {
		return err
	}

	if err := validateConcurrency(c.Concurrency); err != nil {
		return err
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("error: timeout must be greater than 0")
	}

	if c.Requests < 0 {
		return fmt.Errorf("error: requests cannot be negative")
	}

	// Must pick EITHER requests (-n) OR duration (-d)
	if c.Requests > 0 && c.Duration > 0 {
		return fmt.Errorf("error: cannot specify both total requests (-n) and duration (-d)")
	}

	// If neither requests not duration set, return an error
	if c.Requests <= 0 && c.Duration <= 0 {
		return fmt.Errorf("error: must specify either total requests (-n) or duration (-d)")
	}

	return nil
}

func validateConcurrency(concurrency int) error {
	if concurrency < 1 {
		return fmt.Errorf("error: concurrency level must be greater than 0")
	}
	return nil
}

func validateURL(URL string) error {
	URL = strings.TrimSpace(URL)
	if URL == "" {
		return fmt.Errorf("error: target URL (-u / --url) is required")
	}

	requestURL, err := url.ParseRequestURI(URL)
	if err != nil {
		return err
	}

	if requestURL.Scheme != "http" && requestURL.Scheme != "https" {
		return fmt.Errorf("error: target URL scheme must be http or https")
	}

	if requestURL.Host == "" {
		return fmt.Errorf("error: target URL must have a hostname")
	}

	return nil
}

func validateHttpMethod(method string) error {
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return fmt.Errorf("error: invalid HTTP method: %s", method)
	}
	return nil
}
