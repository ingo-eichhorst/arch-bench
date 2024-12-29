package main

import (
	"fmt"
	"os"

	"github.com/ingo-eichhorst/arch-bench/internal/adapters/cli"
	"github.com/ingo-eichhorst/arch-bench/internal/adapters/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	rootCmd := cli.NewRootCmd(cfg)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
