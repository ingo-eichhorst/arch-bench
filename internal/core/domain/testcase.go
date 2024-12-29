package domain

type TestCaseConfig struct {
	Name     string
	Input    string
	Expected string
}

type TestCase struct {
	Name     string
	Input    string
	Expected string
	Result   *TestResult
}

type TestResult struct {
	Output  string
	Metrics []Metric
}
