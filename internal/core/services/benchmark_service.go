package services

import (
	"fmt"
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

	// for each cfg.TeastSuiteConfigs start RunTestSuite sequential
	for _, testSuiteConfig := range s.cfg.TestSuiteConfigs {
		result, err := s.RunTestSuite(testSuiteConfig)
		if err != nil {
			return fmt.Errorf("error running test suite %s: %v", testSuiteConfig.Name, err)
		}
		fmt.Printf("Result: %v\n", result)
	}

	return nil
}

func (s *BenchmarkService) RunTestSuite(cfg domain.TestSuiteConfig) (result *domain.TestSuite, error error) {

	// testSuite is a slice of type TestSuite
	testSuite := domain.TestSuite{
		TestCases: make([]domain.TestCase, len(cfg.TestCaseConfigs)),
	}

	for i, testCaseConfig := range cfg.TestCaseConfigs {
		// create a testcase from the testcase config
		var testCase domain.TestCase

		testCase, err := s.RunTestCase(&cfg, &testCaseConfig)
		if err != nil {
			return nil, fmt.Errorf("error running test case %s: %v", testCaseConfig.Name, err)
		}

		testSuite.TestCases[i] = testCase

	}

	// TODO: Save the results to the repository
	// TODO: Generate and output the CLI report

	return &testSuite, nil
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
	response, err := llmService.GenerateResponse(testCaseConfig.Input, testCaseConfig.Input)
	if err != nil {
		return domain.TestCase{}, fmt.Errorf("error creating LLM response: %v", err)
	}
	endTime := time.Now()

	// for every test case check the metrics for the response
	metrics := s.metricService.CalculateMetrics(
		response,
		testCaseConfig.Expected,
		startTime,
		endTime,
	)

	return domain.TestCase{
		Name:     testCaseConfig.Name,
		Input:    testCaseConfig.Input,
		Expected: testCaseConfig.Expected,
		Result: &domain.TestResult{
			Output:  response,
			Metrics: metrics,
		},
	}, nil
}
