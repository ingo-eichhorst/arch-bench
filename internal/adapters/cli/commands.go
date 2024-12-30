package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/ingo-eichhorst/arch-bench/internal/adapters/config"
	"github.com/ingo-eichhorst/arch-bench/internal/adapters/llm"
	"github.com/ingo-eichhorst/arch-bench/internal/core/ports"
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			benchmarkName := args[0]
			testSuiteName, _ := cmd.Flags().GetString("test-suite")
			return runBenchmark(cfg, benchmarkName, testSuiteName)
		},
	}
	runCmd.Flags().String("test-suite", "", "Specify a test suite to run")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available resources",
	}

	listBenchmarksCmd := &cobra.Command{
		Use:   "benchmarks",
		Short: "List available benchmarks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listBenchmarks()
		},
	}

	listTestSuitesCmd := &cobra.Command{
		Use:   "test-suites <benchmark-name>",
		Short: "List available test suites for a benchmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTestSuites(args[0])
		},
	}

	listProvidersCmd := &cobra.Command{
		Use:   "providers",
		Short: "List available providers and models",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProviders()
		},
	}

	listCmd.AddCommand(listBenchmarksCmd, listTestSuitesCmd, listProvidersCmd)
	rootCmd.AddCommand(runCmd, listCmd)
	return rootCmd
}

func runBenchmark(cfg *config.Config, benchmarkName, testSuiteName string) error {
	benchConfigLoader, err := config.NewBenchmarkConfigLoader(benchmarkName)
	if err != nil {
		return fmt.Errorf("error initiating the benchmark config loader: %v", err)
	}
	benchConfig, err := benchConfigLoader.LoadBenchmarkConfig(
		cfg.EvalAPIKey,
		cfg.EvalModel,
		cfg.EvalProvider,
		cfg.OpenAIAPIKey,
	)
	if err != nil {
		return fmt.Errorf("error loading benchmark config: %v", err)
	}
	service := services.NewBenchmarkService(benchConfig)
	return service.RunBenchmark(testSuiteName)
}

func listBenchmarks() error {
	benchmarksDir := filepath.Join("../", "../", "benchmarks")
	entries, err := os.ReadDir(benchmarksDir)
	if err != nil {
		return fmt.Errorf("error reading benchmarks directory: %v", err)
	}

	fmt.Println("Available benchmarks:")
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println("-", entry.Name())
		}
	}
	return nil
}

func listTestSuites(benchmarkName string) error {
	benchmarkDir := filepath.Join("../", "../", "benchmarks", benchmarkName)
	entries, err := os.ReadDir(benchmarkDir)
	if err != nil {
		return fmt.Errorf("error reading benchmark directory: %v", err)
	}

	fmt.Printf("Test suites for benchmark '%s':\n", benchmarkName)
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println("-", entry.Name())
		}
	}
	return nil
}

func listProviders() error {
	providerDir := filepath.Join("../", "../", "internal", "adapters", "llm")
	entries, err := os.ReadDir(providerDir)
	if err != nil {
		return fmt.Errorf("error reading providers directory: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight) // Align right
	fmt.Fprintf(w, "%-10s\t%-20s\n", "Provider", "Model")                   //Fixed width for better alignment
	fmt.Fprintf(w, "%-10s\t%-20s\n", "--------", "-----")

	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			providerName := strings.TrimSuffix(fileName, ".go")

			var provider ports.LLMProvider
			switch providerName {
			case "openai":
				provider = llm.NewOpenAIProvider("dummy-key", "dummy-model")
			// Add cases for other providers as they are implemented
			default:
				continue
			}

			for _, model := range provider.GetModels() {
				fmt.Fprintf(w, "%-10s\t%-20s\n", providerName, model) //Fixed width
			}
		}
	}

	w.Flush()
	return nil
}

func Execute(cfg *config.Config) error {
	rootCmd := NewRootCmd(cfg)
	return rootCmd.Execute()
}
