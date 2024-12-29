package services

import (
	"fmt"

	"github.com/ingo-eichhorst/arch-bench/internal/adapters/llm"
	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
	"github.com/ingo-eichhorst/arch-bench/internal/core/ports"
)

type LLMService struct {
	provider ports.LLMProvider
}

func NewLLMService(providerName string, APIKey string, ModelName string) (*LLMService, error) {
	var provider ports.LLMProvider
	switch providerName {
	case "openai":
		provider = llm.NewOpenAIProvider(APIKey, ModelName)
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", providerName)
	}

	return &LLMService{provider: provider}, nil
}

func (s *LLMService) GenerateResponse(SystemPrompt string, query string) (domain.LLMResponse, error) {
	llmResponse, err := s.provider.GenerateResponse(SystemPrompt, query)
	if err != nil {
		return domain.LLMResponse{}, err
	}

	return llmResponse, nil
}
