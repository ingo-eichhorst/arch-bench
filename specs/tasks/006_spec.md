# Task #6 - Structured Outputs

Add structured outputs to the OpenAI adapter.

## Objectives
- There is a method that expects the query and a json-schema to provide structured output.
- The method should return an object with the properties defined in the schema.
- The schema will be used in metrics_service.go to call the LLM for validation expecting the response score as a set of schema values.