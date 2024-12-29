package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
)

type BenchmarkService struct {
	cfg           *domain.BenchmarkConfig
	metricService *MetricService
}

func NewBenchmarkService(benchConfig *domain.BenchmarkConfig) *BenchmarkService {
	return &BenchmarkService{
		cfg: benchConfig,
		metricService: NewMetricService(
			benchConfig.EvalProvider,
			benchConfig.EvalModel,
			benchConfig.EvalApiKey,
		),
	}
}

func (s *BenchmarkService) RunBenchmark() error {
	fmt.Printf("Running benchmark %s\n", s.cfg.Name)

	for _, testSuiteConfig := range s.cfg.TestSuiteConfigs {
		testSuite, err := s.RunTestSuite(testSuiteConfig)
		if err != nil {
			return fmt.Errorf("error running test suite %s: %v", testSuiteConfig.Name, err)
		}

		// Generate and output the CLI report for this test suite
		s.outputTestSuiteResults(testSuite)
	}

	return nil
}

func (s *BenchmarkService) RunTestSuite(cfg domain.TestSuiteConfig) (*domain.TestSuite, error) {
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

	return testSuite, nil
}

func (s *BenchmarkService) RunTestCase(testSuiteConfig *domain.TestSuiteConfig, testCaseConfig *domain.TestCaseConfig) (domain.TestCase, error) {
	llmService, err := NewLLMService(
		s.cfg.EvalProvider,
		s.cfg.EvalApiKey,
		s.cfg.EvalModel,
	)
	if err != nil {
		return domain.TestCase{}, fmt.Errorf("error creating LLM service: %v", err)
	}

	startTime := time.Now()
	llmResponse, err := llmService.GenerateResponse(testCaseConfig.Input, testCaseConfig.Input)
	if err != nil {
		return domain.TestCase{}, fmt.Errorf("error creating LLM response: %v", err)
	}
	duration := time.Since(startTime)

	metrics := s.metricService.CalculateMetrics(
		llmResponse.Response,
		testCaseConfig.Expected,
	)

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

func (s *BenchmarkService) outputTestSuiteResults(testSuite *domain.TestSuite) {
	results := testSuite.AggregateResults()

	fmt.Printf("\nResults for Test Suite: %s\n", testSuite.Name)
	fmt.Printf("%-20s %-20s %-15s %-10s %-10s\n", "TestSuite", "TestCase", "Duration", "Cost", "Rating")
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range results {
		fmt.Printf("%-20s %-20s %-15s $%-9.2f %-10.2f\n",
			result.TestSuite,
			result.TestCase,
			result.Duration.Round(time.Millisecond),
			result.Cost,
			result.AverageRating,
		)
	}
	fmt.Println()
}
