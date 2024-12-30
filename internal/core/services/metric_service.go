package services

import (
	"fmt"
	"strconv"

	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
)

type MetricService struct {
	llmService *LLMService
}

func NewMetricService(
	EvalProvider string,
	EvalModel string,
	cfg *domain.BenchmarkConfig,
) *MetricService {
	llmService, _ := NewLLMService(EvalProvider, EvalModel, cfg)
	return &MetricService{
		llmService: llmService,
	}
}

// GEval represents the main evaluator structure
type GEval struct {
	llmService      *LLMService
	TaskPrompt      string
	EvalCriteria    string
	ChainOfThoughts string
}

// NewGEval creates a new G-Eval instance
func NewGEval(llmService *LLMService, taskPrompt, evalCriteria string) *GEval {
	return &GEval{
		llmService:   llmService,
		TaskPrompt:   taskPrompt,
		EvalCriteria: evalCriteria,
	}
}

// GenerateChainOfThoughts generates the evaluation steps
func (g *GEval) GenerateChainOfThoughts() error {
	prompt := fmt.Sprintf("Given the task: %s\nAnd the evaluation criteria: %s\nGenerate a step-by-step chain of thoughts for evaluation:", g.TaskPrompt, g.EvalCriteria)
	response, err := g.llmService.GenerateResponse(prompt, "")
	if err != nil {
		return fmt.Errorf("failed to generate chain of thoughts: %v", err)
	}
	g.ChainOfThoughts = response.Response
	return nil
}

// buildPrompt constructs the full prompt for evaluation
func (g *GEval) buildPrompt(context, target string) string {
	return fmt.Sprintf(`%s
		Evaluation Criteria:
		%s
		Evaluation Steps:
		%s
		Input Context:
		%s
		Input Target:
		%s
		Evaluation Form (scores ONLY):`,
		g.TaskPrompt,
		g.EvalCriteria,
		g.ChainOfThoughts,
		context,
		target,
	)
}

// Evaluate performs the evaluation of the target text given the context
func (g *GEval) Evaluate(context, target string, possibleScores []int) (*domain.EvaluationResult, error) {
	prompt := g.buildPrompt(context, target)
	// TODO: Use Structured Outputs instead to make everything more simple to parse
	response, err := g.llmService.GenerateResponse("You will evaluate the quality of a generated text.", prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get evaluation: %v", err)
	}

	pos := len(response.Response) - 1
	character := response.Response[pos]
	score, _ := strconv.Atoi(string(character))

	return &domain.EvaluationResult{
		Score: float64(score * 20),
	}, nil
}

func (s *MetricService) CalculateMetrics(response string, expected string) []domain.Metric {
	geval := NewGEval(s.llmService, "Evaluate the quality of the generated text.", "Coherence (1-5): evaluate the logical flow and connection between sentences.")

	// TODO: Proper error handling
	err := geval.GenerateChainOfThoughts()
	if err != nil {
		// Handle error, perhaps log it
		fmt.Printf("Error generating chain of thoughts: %v\n", err)
	}

	result, err := geval.Evaluate(expected, response, []int{1, 2, 3, 4, 5})
	if err != nil {
		// Handle error, perhaps log it
		fmt.Printf("Error evaluating response: %v\n", err)
	}

	metrics := []domain.Metric{
		{
			Name:  "geval",
			Value: result.Score,
		},
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
