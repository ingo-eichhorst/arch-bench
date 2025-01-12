package domain

type Benchmark struct {
	Name         string
	EvalProvider string
	EvalModel    string
	TestSuites   []*TestSuite
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

type BenchmarkConfig struct {
	Name             string
	Description      string
	Version          string
	EvalApiKey       string
	EvalModel        string
	EvalProvider     string
	OpenAIAPIKey     string
	TestSuiteConfigs []TestSuiteConfig
}

type EvaluationResult struct {
	Score float64
}
