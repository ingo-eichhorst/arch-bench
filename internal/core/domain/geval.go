package domain

import "encoding/json"

type StructuredOutput struct {
	Type                 string                     `json:"type"`
	Properties           StructuredOutputProperties `json:"properties"`
	Required             []string                   `json:"required"`
	AdditionalProperties bool                       `json:"additionalProperties"`
}

type StructuredOutputProperties struct {
	Score StructuredOutputPropertiy `json:"score"`
}

type StructuredOutputPropertiy struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	// Minimum     float64 `json:"minimum"` // Not yet supported by OpenAI
	// Maximum     float64 `json:"maximum"`
}

func (g StructuredOutput) MarshalJSON() ([]byte, error) {
	type Alias StructuredOutput
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(g),
	})
}
