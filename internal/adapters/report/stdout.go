package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
	"github.com/ingo-eichhorst/arch-bench/internal/core/ports"
)

type StdoutReportGenerator struct{}

func NewStdoutReportCreator() ports.ReportCreator {
	return &StdoutReportGenerator{}
}

func (s *StdoutReportGenerator) GenerateTestSuiteReport(testSuite *domain.TestSuite) error {
	results := testSuite.AggregateResults()

	fmt.Printf("\nResults for Test Suite: %s\n", testSuite.Name)
	fmt.Printf("%-20s %-20s %-15s %-15s %-10s\n", "TestSuite", "TestCase", "Duration", "Cost", "Rating") //Increased width for cost
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range results {
		fmt.Printf("%-20s %-20s %-15s $%-15.6f %-10.2f\n", // Changed to %-15.6f
			result.TestSuite,
			result.TestCase,
			result.Duration.Round(time.Millisecond),
			result.Cost,
			result.AverageRating,
		)
	}
	fmt.Println()
	return nil
}

func (s *StdoutReportGenerator) GenerateBenchmarkReport(benchmark *domain.Benchmark, testSuites []*domain.TestSuite) {
	fmt.Printf("\nBenchmark Results: %s\n", benchmark.Name)
	fmt.Printf("%-20s %-15s %-15s %-10s\n", "TestSuite", "Duration", "Cost", "Avg Rating")
	fmt.Println(strings.Repeat("-", 60))

	var totalBenchmarkDuration time.Duration
	var totalBenchmarkCost float64
	var totalBenchmarkRating float64
	var numTestCases int

	for _, testSuite := range testSuites {
		results := testSuite.AggregateResults()
		var totalTestSuiteDuration time.Duration
		var totalTestSuiteCost float64
		var totalTestSuiteRating float64

		for _, result := range results {
			totalTestSuiteDuration += result.Duration
			totalTestSuiteCost += result.Cost
			totalTestSuiteRating += result.AverageRating
			numTestCases++
		}

		avgRating := totalTestSuiteRating / float64(len(results))

		fmt.Printf("%-20s %-15s $%-15.6f %-10.2f\n",
			testSuite.Name,
			totalTestSuiteDuration.Round(time.Millisecond),
			totalTestSuiteCost,
			avgRating,
		)

		totalBenchmarkDuration += totalTestSuiteDuration
		totalBenchmarkCost += totalTestSuiteCost
		totalBenchmarkRating += totalTestSuiteRating
	}

	avgBenchmarkRating := totalBenchmarkRating / float64(numTestCases)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Benchmark Summary:\n")
	fmt.Printf("Total Duration: %-15s\n", totalBenchmarkDuration.Round(time.Millisecond))
	fmt.Printf("Total Cost:     $%-15.6f\n", totalBenchmarkCost)
	fmt.Printf("Average Rating: %-10.2f\n", avgBenchmarkRating)
	fmt.Println()
}
