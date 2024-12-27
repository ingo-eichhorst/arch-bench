package ports

import "github.com/ingo-eichhorst/arch-bench/internal/core/domain"

type BenchmarkRunner interface {
	Run(benchmark *domain.Benchmark) error
}
