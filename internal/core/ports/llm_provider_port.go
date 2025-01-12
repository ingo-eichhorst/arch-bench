package ports

import "github.com/ingo-eichhorst/arch-bench/internal/core/domain"

type LLMProvider interface {
	GenerateResponse(systemPrompt string, query string, images []string) (domain.LLMResponse, error)
	GetModels() []string
	GenerateStructuredResponse(systemPrompt string, query string, schema domain.StructuredOutput) (map[string]interface{}, error)
}
