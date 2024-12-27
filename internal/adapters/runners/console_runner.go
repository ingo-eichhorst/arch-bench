package runners

import (
	"fmt"

	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
)

type ConsoleBenchmarkRunner struct{}

func (r *ConsoleBenchmarkRunner) Run(benchmark *domain.Benchmark) error {
	fmt.Printf("Running benchmark: %s\n", benchmark.Name)
	fmt.Println("Demo Result: Success!")
	return nil
}
