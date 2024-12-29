package cli

import (
	"fmt"

	"github.com/ingo-eichhorst/arch-bench/internal/adapters/config"
	"github.com/ingo-eichhorst/arch-bench/internal/core/services"
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "arch-bench",
		Short: "CLI for software architecture benchmark for LLMs and LVMs",
	}

	runCmd := &cobra.Command{
		Use:   "run <benchmark-name>",
		Short: "Run a benchmark",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			benchmarkName := args[0]
			benchConfigLoader, err := config.NewBenchmarkConfigLoader(benchmarkName)
			if err != nil {
				return fmt.Errorf("error initiating the benchmark config loader: %v", err)
			}
			benchConfig, err := benchConfigLoader.LoadBenchmarkConfig(
				cfg.EvalAPIKey,
				cfg.EvalModel,
				cfg.EvalProvider,
			)
			if err != nil {
				return fmt.Errorf("error loading benchmark config: %v", err)
			}
			service := services.NewBenchmarkService(benchConfig)
			return service.RunBenchmark()
		},
	}

	rootCmd.AddCommand(runCmd)
	return rootCmd
}

func Execute(cfg *config.Config) error {
	rootCmd := NewRootCmd(cfg)
	return rootCmd.Execute()
}
