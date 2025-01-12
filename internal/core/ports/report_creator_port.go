package ports

import "github.com/ingo-eichhorst/arch-bench/internal/core/domain"

type ReportCreator interface {
	GenerateTestSuiteReport(testSuite *domain.TestSuite) error
	GenerateBenchmarkReport(benchmark *domain.Benchmark, testSuites []*domain.TestSuite)
}
