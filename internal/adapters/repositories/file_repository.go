package repositories

import "github.com/ingo-eichhorst/arch-bench/internal/core/domain"

type FileBenchmarkRepository struct{}

func (r *FileBenchmarkRepository) GetBenchmark(name string) (*domain.Benchmark, error) {
	// For demo purposes, always return a demo benchmark
	return &domain.Benchmark{Name: "Demo Benchmark"}, nil
}
