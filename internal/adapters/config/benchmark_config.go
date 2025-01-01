package config

import (
	"encoding/json"
	"fmt"
	"io"
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
	OpenAIAPIKey string,
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
	benchmarkConfig.OpenAIAPIKey = OpenAIAPIKey

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

	// load and parse the config.json in the test case folder
	configPath := filepath.Join(casePath, "config.json")
	configFile, err := os.Open(configPath)
	if err != nil {
		return domain.TestCaseConfig{}, fmt.Errorf("failed to open config.json: %w", err)
	}
	defer configFile.Close()

	configData, err := io.ReadAll(configFile)
	if err != nil {
		return domain.TestCaseConfig{}, fmt.Errorf("failed to read config.json: %w", err)
	}

	var config = domain.TestCaseConfig{}
	if err := json.Unmarshal(configData, &config); err != nil {
		return domain.TestCaseConfig{}, fmt.Errorf("failed to parse config.json: %w", err)
	}

	// Construct full paths based on config values, not hardcoded filenames
	inputPath := filepath.Join(casePath, config.Input)
	inputBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return domain.TestCaseConfig{}, fmt.Errorf("failed to read input file '%s': %w", inputPath, err)
	}

	expectedOutputPath := filepath.Join(casePath, config.Expected)
	expectedOutputBytes, err := ioutil.ReadFile(expectedOutputPath)
	if err != nil {
		return domain.TestCaseConfig{}, fmt.Errorf("failed to read expected output file '%s': %w", expectedOutputPath, err)
	}

	config = domain.TestCaseConfig{
		Name:     caseName,
		Path:     casePath,
		Input:    string(inputBytes),
		Expected: string(expectedOutputBytes),
		Images:   config.Images,
	}

	return config, nil
}
