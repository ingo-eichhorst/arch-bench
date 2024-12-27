package services

import (
	"github.com/ingo-eichhorst/arch-bench/internal/core/ports"
)

type BenchmarkService struct {
	repo   ports.BenchmarkRepository
	runner ports.BenchmarkRunner
}

func NewBenchmarkService(repo ports.BenchmarkRepository, runner ports.BenchmarkRunner) *BenchmarkService {
	return &BenchmarkService{repo: repo, runner: runner}
}

func (s *BenchmarkService) RunBenchmark(name string) error {
	benchmark, err := s.repo.GetBenchmark(name)
	if err != nil {
		return err
	}
	return s.runner.Run(benchmark)
}
