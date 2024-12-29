package domain

type LLMResponse struct {
	Response string  `json:"response"`
	Cost     float64 `json:"cost"`
}
