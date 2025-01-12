package llm

import (
	"context"
	"encoding/json"
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
	"gpt-4":                {InputCostPerMillion: 30.0, OutputCostPerMillion: 60.0},
	"gpt-4-32k":            {InputCostPerMillion: 60.0, OutputCostPerMillion: 120.0},
	"gpt-4-turbo":          {InputCostPerMillion: 10.0, OutputCostPerMillion: 30.0},
	"gpt-4-vision-preview": {InputCostPerMillion: 10.0, OutputCostPerMillion: 30.0},
	"gpt-4o":               {InputCostPerMillion: 2.5, OutputCostPerMillion: 10.0},
	"gpt-4o-mini":          {InputCostPerMillion: 0.15, OutputCostPerMillion: 0.6},
	"o1":                   {InputCostPerMillion: 15.0, OutputCostPerMillion: 60.0},
	"o1-preview":           {InputCostPerMillion: 15.0, OutputCostPerMillion: 60.0},
	"o1-mini":              {InputCostPerMillion: 3.0, OutputCostPerMillion: 12.0},
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

func (p *OpenAIProvider) GenerateResponse(systemPrompt string, query string, images []string) (domain.LLMResponse, error) {
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

func (p *OpenAIProvider) GetModels() []string {
	models := make([]string, 0, len(ModelPriceMap))
	for model := range ModelPriceMap {
		models = append(models, model)
	}
	return models
}

// GenerateStructuredResponse generates a response with structured output based on a JSON schema.
func (p *OpenAIProvider) GenerateStructuredResponse(systemPrompt, query string, schema domain.StructuredOutput) (map[string]interface{}, error) {

	// convert schema jsonstring to openai.ChatCompletionResponseFormat
	responseFormat := openai.ChatCompletionResponseFormat{
		Type: "json_schema", // This is always "json_object" for JSON schema
		JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
			Name:        "GEvalEvaluationScore",
			Description: "Reasons about the evaluation score or a LLM or LVM completion",
			Schema:      schema,
			Strict:      true,
		},
	}

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
			ResponseFormat: &responseFormat,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error calling OpenAI API: %v", err)
	}

	// Parse the JSON response.  Error handling is crucial here.
	var data map[string]interface{}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	cost, err := p.costCalculator.CalculateCost(resp.Usage.PromptTokens, resp.Usage.CompletionTokens, p.model)
	if err != nil {
		return nil, fmt.Errorf("error calculating cost: %v", err)
	}

	data["cost"] = cost
	return data, nil
}
