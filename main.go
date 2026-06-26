package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mohyehia/goku/internal/config"
	"github.com/mohyehia/goku/internal/engine"
	"github.com/mohyehia/goku/internal/metrics"
)

func main() {
	cfg, err := config.ParseFlags()

	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("Usage Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("Usage Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Start Execution Engine
	fmt.Println("==================================================")
	fmt.Printf("🔥 GOKU LOAD TEST ENGINE STARTING 🔥\n")
	fmt.Println("==================================================")
	fmt.Printf("🎯 Target URL:   %s\n", cfg.URL)
	fmt.Printf("⚡ HTTP Method:   %s\n", cfg.Method)
	fmt.Printf("👥 Concurrency:  %d workers\n", cfg.Concurrency)
	fmt.Printf("⏱️ Request TO:   %v\n", cfg.Timeout)

	// Print headers if exist
	if len(cfg.Headers) > 0 {
		fmt.Println("📋 Custom Headers:")
		for k, v := range cfg.Headers {
			fmt.Printf("   ├── %s: %s\n", k, v)
		}
	}
	if cfg.Requests > 0 {
		fmt.Printf("📊 Profile:      Count-based (%d requests)\n", cfg.Requests)
	} else {
		fmt.Printf("📊 Profile:      Duration-based (%v limit)\n", cfg.Duration)
	}

	if len(cfg.RequestPayload) > 0 {
		fmt.Printf("Payload Length:  %d\n", len(cfg.RequestPayload))
	}

	aggregator := metrics.NewAggregator()

	engine.Orchestrate(cfg, aggregator)

	aggregator.Aggregate()

	fmt.Println("🏁 Pipeline finished!")
}
