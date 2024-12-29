package llm

import (
	"context"
	"fmt"

	"github.com/ingo-eichhorst/arch-bench/internal/core/domain"
	"github.com/ingo-eichhorst/arch-bench/internal/core/ports"
	"github.com/sashabaranov/go-openai"
)

type PriceEntry struct {
	InputCostPerMillion  float64
	OutputCostPerMillion float64
}

var ModelPriceMap = map[string]PriceEntry{
	"gpt-4":                  {InputCostPerMillion: 30.0, OutputCostPerMillion: 60.0},
	"gpt-4-32k":              {InputCostPerMillion: 60.0, OutputCostPerMillion: 120.0},
	"gpt-4-turbo":            {InputCostPerMillion: 10.0, OutputCostPerMillion: 30.0},
	"gpt-4-turbo-2024-04-09": {InputCostPerMillion: 10.0, OutputCostPerMillion: 30.0},
	"gpt-4-0125-preview":     {InputCostPerMillion: 10.0, OutputCostPerMillion: 30.0},
	"gpt-4-1106-preview":     {InputCostPerMillion: 10.0, OutputCostPerMillion: 30.0},
	"gpt-4-vision-preview":   {InputCostPerMillion: 10.0, OutputCostPerMillion: 30.0},
	"gpt-4o":                 {InputCostPerMillion: 2.5, OutputCostPerMillion: 10.0},
	"gpt-4o-2024-11-20":      {InputCostPerMillion: 2.5, OutputCostPerMillion: 10.0},
	"gpt-4o-2024-08-06":      {InputCostPerMillion: 2.5, OutputCostPerMillion: 10.0},
	"gpt-4o-2024-05-13":      {InputCostPerMillion: 5.0, OutputCostPerMillion: 15.0},
	"gpt-4o-mini":            {InputCostPerMillion: 0.15, OutputCostPerMillion: 0.6},
	"gpt-4o-mini-2024-07-18": {InputCostPerMillion: 0.15, OutputCostPerMillion: 0.6},
	"o1":                     {InputCostPerMillion: 15.0, OutputCostPerMillion: 60.0},
	"o1-2024-12-17":          {InputCostPerMillion: 15.0, OutputCostPerMillion: 60.0},
	"o1-preview":             {InputCostPerMillion: 15.0, OutputCostPerMillion: 60.0},
	"o1-preview-2024-09-12":  {InputCostPerMillion: 15.0, OutputCostPerMillion: 60.0},
	"o1-mini":                {InputCostPerMillion: 3.0, OutputCostPerMillion: 12.0},
	"o1-mini-2024-09-12":     {InputCostPerMillion: 3.0, OutputCostPerMillion: 12.0},
}

type OpenAIProvider struct {
	client         *openai.Client
	model          string
	costCalculator CostCalculator
}

// CostCalculator interface defines the contract for calculating costs.  This allows for easier mocking and swapping of cost calculation logic.
type CostCalculator interface {
	CalculateCost(promptTokens int, completionTokens int, model string) (float64, error)
}

type MapBasedCostCalculator struct{}

func (c *MapBasedCostCalculator) CalculateCost(promptTokens, completionTokens int, model string) (float64, error) {
	price, ok := ModelPriceMap[model]
	if !ok {
		return 0, fmt.Errorf("model not found in price map: %s", model)
	}

	promptCost := (float64(promptTokens) / 1000000) * price.InputCostPerMillion
	completionCost := (float64(completionTokens) / 1000000) * price.OutputCostPerMillion

	return promptCost + completionCost, nil
}

func NewOpenAIProvider(apiKey, model string) ports.LLMProvider {
	client := openai.NewClient(apiKey)
	return &OpenAIProvider{
		client:         client,
		model:          model,
		costCalculator: &MapBasedCostCalculator{},
	}
}

func (p *OpenAIProvider) GenerateResponse(systemPrompt string, query string) (domain.LLMResponse, error) {
	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)

	if err != nil {
		return domain.LLMResponse{}, fmt.Errorf("error calling OpenAI API: %v", err)
	}

	cost, err := p.costCalculator.CalculateCost(resp.Usage.PromptTokens, resp.Usage.CompletionTokens, p.model)
	if err != nil {
		return domain.LLMResponse{}, fmt.Errorf("error calculating cost: %v", err)
	}

	return domain.LLMResponse{
		Response: resp.Choices[0].Message.Content,
		Cost:     cost,
	}, nil
}
