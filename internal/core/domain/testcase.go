package domain

import "time"

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
	Output   string
	Metrics  []Metric
	Duration time.Duration
	Cost     float64
}

func (tc *TestCase) CalculateAverageRating() float64 {
	if tc.Result == nil || len(tc.Result.Metrics) == 0 {
		return 0
	}

	var sum float64
	var count int

	for _, metric := range tc.Result.Metrics {
		if metric.Name != "duration" && metric.Name != "cost" {
			sum += metric.Value
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}
