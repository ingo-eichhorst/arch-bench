package domain

type Benchmark struct {
	Name       string
	TestSuites []*TestSuite
}

type MeasurementConfig struct {
	Name   string
	Unit   string
	Weight float64
}

type MetricConfig struct {
	Name        string
	Mesurements []MeasurementConfig
}

type TestSuiteConfig struct {
	Name            string
	Description     string
	MetricConfigs   []MetricConfig
	TestCaseConfigs []TestCaseConfig
}

type BenchmarkConfig struct {
	Name             string
	Description      string
	Version          string
	EvalApiKey       string
	EvalModel        string
	EvalProvider     string
	TestSuiteConfigs []TestSuiteConfig
}

type EvaluationResult struct {
	Score float64
}
