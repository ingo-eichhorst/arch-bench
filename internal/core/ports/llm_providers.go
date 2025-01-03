package ports

import "github.com/ingo-eichhorst/arch-bench/internal/core/domain"

type LLMProvider interface {
	GenerateResponse(SystemPrompt string, query string) (domain.LLMResponse, error)
	GetModels() []string
	GenerateStructuredResponse(systemPrompt, query, schema string) (map[string]interface{}, error)
}
