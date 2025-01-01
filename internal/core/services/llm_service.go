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

func NewLLMService(providerName string, ModelName string, cfg *domain.BenchmarkConfig) (*LLMService, error) {
	var provider ports.LLMProvider
	switch providerName {
	case "openai":
		APIKey := cfg.OpenAIAPIKey
		if APIKey == "" {
			return nil, fmt.Errorf("OpenAI API key is required")
		}
		provider = llm.NewOpenAIProvider(APIKey, ModelName)
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", providerName)
	}

	return &LLMService{provider: provider}, nil
}

func (s *LLMService) GenerateResponse(SystemPrompt string, query string) (domain.LLMResponse, error) {
	// TODO: FAKE images for now
	images := make([]string, 0)
	llmResponse, err := s.provider.GenerateResponse(SystemPrompt, query, images)
	if err != nil {
		return domain.LLMResponse{}, err
	}

	return llmResponse, nil
}

// Add the new GetModelPriceMap method
func (s *LLMService) GetModels() []string {
	return s.provider.GetModels()
}
