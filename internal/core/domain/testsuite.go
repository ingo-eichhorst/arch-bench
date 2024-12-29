package domain

type TestSuite struct {
	Name      string
	Metrics   []Metric
	TestCases []TestCase
}

type Metric struct {
	Name  string
	Value float64
}

type Provider struct {
	Name  string
	Model string
}
