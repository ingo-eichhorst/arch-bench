package services

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

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

func (s *LLMService) GenerateResponse(systemPrompt string, query string, images []string) (domain.LLMResponse, error) {
	// Generate image embeddings (using a separate vision model)
	base64Images := make([]string, len(images))
	for i, imagePath := range images {
		base64Image, err := encodeImageToBase64(imagePath)
		if err != nil {
			return domain.LLMResponse{}, fmt.Errorf("error encoding image to base64: %w", err)
		}
		base64Images[i] = base64Image
	}

	llmResponse, err := s.provider.GenerateResponse(systemPrompt, query, base64Images)
	if err != nil {
		return domain.LLMResponse{}, fmt.Errorf("error calling LLM provider: %w", err)
	}
	return llmResponse, nil
}

func encodeImageToBase64(imagePath string) (string, error) {
	// Read the image file
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("error reading image file: %w", err)
	}

	// Encode the image to base64
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	return encodedImage, nil
}

func (s *LLMService) GetModels() []string {
	return s.provider.GetModels()
}
