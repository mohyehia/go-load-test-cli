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
	if cfg.Requests > 0 {
		fmt.Printf("📊 Profile:      Count-based (%d requests)\n", cfg.Requests)
	} else {
		fmt.Printf("📊 Profile:      Duration-based (%v limit)\n", cfg.Duration)
	}

	aggregator := metrics.NewAggregator()

	engine.Orchestrate(cfg, aggregator)

	aggregator.Aggregate()

	fmt.Println("🏁 Pipeline finished!")
}
