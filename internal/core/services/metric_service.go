package services

import (
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

func (s *MetricService) CalculateMetrics(response string, expected string) []domain.Metric {
	metrics := []domain.Metric{
		{
			Name:  "relevance",
			Value: calculateRelevance(expected, response),
		},
	}

	return metrics
}

func calculateRelevance(expected, actual string) float64 {
	// TODO: Implement actual relevance calculation
	// This could involve more sophisticated NLP techniques or comparison algorithms
	if expected == actual {
		return 100.0
	}
	return 50.0 // Placeholder implementation
}
