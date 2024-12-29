package ports

type LLMProvider interface {
	GenerateResponse(SystemPrompt string, query string) (string, error)
}
