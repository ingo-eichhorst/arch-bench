package services

import (
	"time"

	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
)

type MetricService struct{}

func NewMetricService(
	EvalProvider string,
	EvalModel string,
	EvalApiKey string,
) *MetricService {
	return &MetricService{}
}

func (s *MetricService) CalculateMetrics(response string, expected string, startTime time.Time, endTime time.Time) []domain.Metric {
	metrics := []domain.Metric{
		{
			Name:  "cost",
			Value: calculateCost(response),
		},
		{
			Name:  "time",
			Value: float64(endTime.Sub(startTime).Milliseconds()),
		},
		{
			Name:  "relevance",
			Value: calculateRelevance(expected, response),
		},
	}

	return metrics
}

func calculateCost(output string) float64 {
	// TODO: Implement actual cost calculation based on the model and output length
	return float64(len(output)) * 0.001 // Placeholder implementation
}

func calculateRelevance(expected, actual string) float64 {
	// TODO: Implement actual relevance calculation
	// This could involve more sophisticated NLP techniques or comparison algorithms
	if expected == actual {
		return 100.0
	}
	return 50.0 // Placeholder implementation
}
