package services

import (
	"fmt"

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
	response, err := g.llmService.GenerateResponse("", prompt, nil)
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

// Evaluate performs the evaluation using structured output.
func (g *GEval) Evaluate(context string, target string) (*domain.EvaluationResult, error) {
	var gevalSchemaVar = domain.StructuredOutput{
		Type: "object",
		Properties: domain.StructuredOutputProperties{
			Score: domain.StructuredOutputPropertiy{
				Type:        "number",
				Description: "Overall evaluation score (0-100)",
				// Minimum:     0, // Not yet supported by OpenAI
				// Maximum:     100,
			},
		},
		Required:             []string{"score"},
		AdditionalProperties: false,
	}

	structuredResponse, err := g.llmService.provider.GenerateStructuredResponse(g.buildPrompt(context, target), target, gevalSchemaVar)
	if err != nil {
		return nil, fmt.Errorf("error generating structured response: %v", err)
	}

	scoreInterface, ok := structuredResponse["score"]
	if !ok {
		return nil, fmt.Errorf("score field not found in structured response")
	}

	return &domain.EvaluationResult{Score: scoreInterface.(float64)}, nil
}

func (s *MetricService) CalculateMetrics(response string, expected string) ([]domain.Metric, error) {
	geval := NewGEval(s.llmService, "Evaluate the quality of the generated text.", "Coherence (1-5): evaluate the logical flow and connection between sentences.")
	err := geval.GenerateChainOfThoughts()
	if err != nil {
		return nil, fmt.Errorf("error generating chain of thoughts: %v", err)
	}

	result, err := geval.Evaluate(expected, response)
	if err != nil {
		return nil, fmt.Errorf("error evaluating response: %v", err)
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

	return metrics, nil
}

func calculateRelevance(expected, actual string) float64 {
	if expected == actual {
		return 100.0
	}
	return 50.0
}
