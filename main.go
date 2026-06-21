package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mohyehia/goku/internal/config"
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
		os.Exit(1)
	}

	// Start Execution Engine
}
