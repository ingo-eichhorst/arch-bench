package ports

import "github.com/ingo-eichhorst/arch-bench/internal/core/domain"

type BenchmarkRepository interface {
	GetBenchmark(name string) (*domain.Benchmark, error)
}
