package services

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ingo-eichhorst/arch-bench/internal/adapters/report"
	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
	"github.com/ingo-eichhorst/arch-bench/internal/core/ports"
)

type BenchmarkService struct {
	cfg           *domain.BenchmarkConfig
	metricService *MetricService
}

type Report struct {
	creator ports.ReportCreator
}

func NewBenchmarkService(benchConfig *domain.BenchmarkConfig) *BenchmarkService {
	return &BenchmarkService{
		cfg: benchConfig,
		metricService: NewMetricService(
			benchConfig.EvalProvider,
			benchConfig.EvalModel,
			benchConfig,
		),
	}
}

func (s *BenchmarkService) RunBenchmark(testSuiteName string) error {
	fmt.Printf("Running Benchmark: %s\n", s.cfg.Name)
	fmt.Printf("- Eval Provider: %s\n", s.cfg.EvalProvider)
	fmt.Printf("- Eval Model: %s\n", s.cfg.EvalModel)
	fmt.Printf("----------------------\n")

	var completedTestSuites []*domain.TestSuite
	var completedBenchmark *domain.Benchmark
	stdOutReport := report.NewStdoutReportCreator()

	for _, testSuiteConfig := range s.cfg.TestSuiteConfigs {
		if testSuiteName != "" && testSuiteConfig.Name != testSuiteName {
			continue // Skip if testSuiteName is specified and doesn't match
		}
		testSuite, err := s.RunTestSuite(testSuiteConfig)
		if err != nil {
			return fmt.Errorf("error running test suite %s: %v", testSuiteConfig.Name, err)
		}
		stdOutReport.GenerateTestSuiteReport(testSuite)
		completedTestSuites = append(completedTestSuites, testSuite)
	}

	completedBenchmark = &domain.Benchmark{
		Name:         s.cfg.Name,
		EvalProvider: s.cfg.EvalProvider,
		EvalModel:    s.cfg.EvalModel,
		TestSuites:   completedTestSuites,
	}
	stdOutReport.GenerateBenchmarkReport(completedBenchmark, completedTestSuites)

	return nil
}

func (s *BenchmarkService) RunTestSuite(cfg domain.TestSuiteConfig) (*domain.TestSuite, error) {
	fmt.Printf("Running Test Suite: %s\n", cfg.Name)

	testSuite := &domain.TestSuite{
		Name:      cfg.Name,
		TestCases: make([]domain.TestCase, len(cfg.TestCaseConfigs)),
	}

	for i, testCaseConfig := range cfg.TestCaseConfigs {
		testCase, err := s.RunTestCase(&cfg, &testCaseConfig)
		if err != nil {
			return nil, fmt.Errorf("error running test case %s: %v", testCaseConfig.Name, err)
		}

		testSuite.TestCases[i] = testCase
	}
	fmt.Printf("----------------------\n")

	return testSuite, nil
}

func (s *BenchmarkService) RunTestCase(testSuiteConfig *domain.TestSuiteConfig, testCaseConfig *domain.TestCaseConfig) (domain.TestCase, error) {
	fmt.Printf("Running Test Case: %s\n", testCaseConfig.Name)

	llmService, err := NewLLMService(
		testSuiteConfig.Provider,
		testSuiteConfig.Model,
		s.cfg,
	)
	if err != nil {
		return domain.TestCase{}, fmt.Errorf("error creating LLM service: %v", err)
	}

	startTime := time.Now()

	// Convert relative image paths to absolute paths
	absImages := make([]string, len(testCaseConfig.Images))
	for i, imagePath := range testCaseConfig.Images {
		absPath, err := filepath.Abs(filepath.Join(testCaseConfig.Path, imagePath))
		if err != nil {
			return domain.TestCase{}, fmt.Errorf("error converting image path to absolute path: %w", err)
		}
		absImages[i] = absPath
	}

	llmResponse, err := llmService.GenerateResponse(testCaseConfig.Input, testCaseConfig.Expected, absImages)
	if err != nil {
		return domain.TestCase{}, fmt.Errorf("error creating LLM response: %v", err)
	}
	duration := time.Since(startTime)

	metrics, err := s.metricService.CalculateMetrics(
		llmResponse.Response,
		testCaseConfig.Expected,
	)
	if err != nil {
		return domain.TestCase{}, fmt.Errorf("error calculating metrics: %v", err)
	}

	return domain.TestCase{
		Name:     testCaseConfig.Name,
		Input:    testCaseConfig.Input,
		Expected: testCaseConfig.Expected,
		Result: &domain.TestResult{
			Output:   llmResponse.Response,
			Metrics:  metrics,
			Duration: duration,
			Cost:     llmResponse.Cost,
		},
	}, nil
}
