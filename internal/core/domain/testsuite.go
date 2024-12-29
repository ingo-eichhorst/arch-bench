package domain

import (
	"time"
)

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

type TestSuiteResult struct {
	TestSuite     string
	TestCase      string
	Duration      time.Duration
	Cost          float64
	AverageRating float64
}

func (ts *TestSuite) AggregateResults() []TestSuiteResult {
	results := make([]TestSuiteResult, len(ts.TestCases))

	for i, testCase := range ts.TestCases {
		results[i] = TestSuiteResult{
			TestSuite:     ts.Name,
			TestCase:      testCase.Name,
			Duration:      testCase.Result.Duration,
			Cost:          testCase.Result.Cost,
			AverageRating: testCase.CalculateAverageRating(),
		}
	}

	return results
}
