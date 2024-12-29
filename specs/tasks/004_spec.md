# Task #4 - CLI Output

Metrik vornehmen und den Output in Tabellenform produzieren. Dann schreibe ich ein readme und

## Objective
  - Add a new field `cost` to the resonse from the LLM. Create a LLMResponseStruct as reference.
  - For example: {"response": "Hello World", "cost": 0.1}

## Docs

To calculate the cost of an OpenAI completion using the github.com/sashabaranov/go-openai library, follow these steps:

Obtain Token Usage: After making a completion request, retrieve the number of tokens used for both the prompt (input) and the completion (output). The OpenAI API typically returns this information in the response's Usage field. In Go, you can access it as follows:

go

Code kopieren

response, err := client.CreateCompletion(ctx, openai.CompletionRequest{ Model: openai.GPT4, Prompt: "Your prompt here", // Additional parameters }) if err != nil { // Handle error } promptTokens := response.Usage.PromptTokens completionTokens := response.Usage.CompletionTokens totalTokens := response.Usage.TotalTokens

Determine Model Pricing: Refer to OpenAI's pricing page to find the cost per 1,000 tokens for the specific model you're using. For example, as of December 2024, the pricing for the gpt-4o model is:

Prompt (input) tokens: $2.50 per 1 million tokens

Completion (output) tokens: $10.00 per 1 million tokens

This translates to:

Prompt tokens: $0.0025 per 1,000 tokens

Completion tokens: $0.01 per 1,000 tokens

Calculate Costs: Compute the cost for both the prompt and completion tokens, then sum them to get the total cost. Here's how you can do it in Go:

go

Code kopieren

// Define the cost per 1,000 tokens in USD const promptCostPerThousand = 0.0025 const completionCostPerThousand = 0.01 // Calculate the cost for prompt and completion promptCost := (float64(promptTokens) / 1000) * promptCostPerThousand completionCost := (float64(completionTokens) / 1000) * completionCostPerThousand // Total cost totalCost := promptCost + completionCost fmt.Printf("Prompt Tokens: %d\n", promptTokens) fmt.Printf("Completion Tokens: %d\n", completionTokens) fmt.Printf("Total Tokens: %d\n", totalTokens) fmt.Printf("Prompt Cost: $%.6f\n", promptCost) fmt.Printf("Completion Cost: $%.6f\n", completionCost) fmt.Printf("Total Cost: $%.6f\n", totalCost)

Note: Ensure you use the latest pricing from OpenAI's official pricing page, as rates may change over time.

By following these steps, you can accurately calculate the cost of each completion request made using the go-openai library.