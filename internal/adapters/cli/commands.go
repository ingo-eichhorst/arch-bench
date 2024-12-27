package cli

import (
	"log"

	"github.com/ingo-eichhorst/arch-bench/internal/adapters/repositories"
	"github.com/ingo-eichhorst/arch-bench/internal/adapters/runners"
	"github.com/ingo-eichhorst/arch-bench/internal/core/services"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark CLI for software architecture tasks",
}

var runCmd = &cobra.Command{
	Use:   "run <benchmark-name>",
	Short: "Run a benchmark",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		benchmarkName := args[0]
		service := services.NewBenchmarkService(&repositories.FileBenchmarkRepository{}, &runners.ConsoleBenchmarkRunner{})
		err := service.RunBenchmark(benchmarkName)
		if err != nil {
			log.Fatalf("Error running benchmark: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
