package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
)

type BenchmarkConfig struct {
	Benchmark domain.Benchmark
}

type BenchmarkConfigLoader struct {
	BasePath string
	Name     string
}

func NewBenchmarkConfigLoader(name string) (*BenchmarkConfigLoader, error) {
	basePath := filepath.Join("../../", "benchmarks", name)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("benchmark directory does not exist: %s", basePath)
	}
	return &BenchmarkConfigLoader{
		BasePath: basePath,
		Name:     name,
	}, nil
}

func (l *BenchmarkConfigLoader) LoadBenchmarkConfig(
	evalAPIKey string,
	evalModel string,
	evalProvider string,
) (*domain.BenchmarkConfig, error) {
	configPath := filepath.Join(l.BasePath, "config.json")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read benchmark config: %w", err)
	}

	var benchmarkConfig domain.BenchmarkConfig
	err = json.Unmarshal(data, &benchmarkConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal benchmark config: %w", err)
	}

	// Merge with default config
	benchmarkConfig.EvalApiKey = evalAPIKey
	benchmarkConfig.EvalModel = evalModel
	benchmarkConfig.EvalProvider = evalProvider

	// Load test suites
	testSuites, err := l.loadTestSuites()
	if err != nil {
		return nil, fmt.Errorf("failed to load test suites: %w", err)
	}
	benchmarkConfig.TestSuiteConfigs = testSuites

	return &benchmarkConfig, nil
}

func (l *BenchmarkConfigLoader) loadTestSuites() ([]domain.TestSuiteConfig, error) {
	var testSuites []domain.TestSuiteConfig

	suiteDirs, err := ioutil.ReadDir(l.BasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read benchmark directory: %w", err)
	}

	for _, suiteDir := range suiteDirs {
		if suiteDir.IsDir() {
			suite, err := l.loadTestSuite(suiteDir.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to load test suite %s: %w", suiteDir.Name(), err)
			}
			testSuites = append(testSuites, suite)
		}
	}

	return testSuites, nil
}

func (l *BenchmarkConfigLoader) loadTestSuite(suiteName string) (domain.TestSuiteConfig, error) {
	suitePath := filepath.Join(l.BasePath, suiteName)
	configPath := filepath.Join(suitePath, "config.json")

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return domain.TestSuiteConfig{}, fmt.Errorf("failed to read test suite config: %w", err)
	}

	var suite domain.TestSuiteConfig
	err = json.Unmarshal(data, &suite)
	if err != nil {
		return domain.TestSuiteConfig{}, fmt.Errorf("failed to unmarshal test suite config: %w", err)
	}

	suite.Name = suiteName

	// Load test cases
	testCases, err := l.loadTestCases(suitePath)
	if err != nil {
		return domain.TestSuiteConfig{}, fmt.Errorf("failed to load test cases: %w", err)
	}
	suite.TestCaseConfigs = testCases

	return suite, nil
}

func (l *BenchmarkConfigLoader) loadTestCases(suitePath string) ([]domain.TestCaseConfig, error) {
	var testCases []domain.TestCaseConfig

	caseDirs, err := ioutil.ReadDir(suitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test suite directory: %w", err)
	}

	for _, caseDir := range caseDirs {
		if caseDir.IsDir() {
			testCase, err := l.loadTestCase(suitePath, caseDir.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to load test case %s: %w", caseDir.Name(), err)
			}
			testCases = append(testCases, testCase)
		}
	}

	return testCases, nil
}

func (l *BenchmarkConfigLoader) loadTestCase(suitePath, caseName string) (domain.TestCaseConfig, error) {
	casePath := filepath.Join(suitePath, caseName)

	inputPath := filepath.Join(casePath, "input.txt")
	input, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return domain.TestCaseConfig{}, fmt.Errorf("failed to read input file: %w", err)
	}

	expectedOutputPath := filepath.Join(casePath, "expected_output.txt")
	expectedOutput, err := ioutil.ReadFile(expectedOutputPath)
	if err != nil {
		return domain.TestCaseConfig{}, fmt.Errorf("failed to read expected output file: %w", err)
	}

	testCase := domain.TestCaseConfig{
		Name:     caseName,
		Input:    string(input),
		Expected: string(expectedOutput),
	}

	return testCase, nil
}
