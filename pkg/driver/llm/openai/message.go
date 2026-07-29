package openai

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// Message is the structured output the model must return.
type Message struct {
	Content   string     `json:"content" jsonschema_description:"Assistant text content for this turn"`
	ToolCalls []ToolCall `json:"tool_calls" jsonschema_description:"Tools to execute; empty when the turn is complete"`
}

// ToolCall is a tool invocation requested by the model.
type ToolCall struct {
	Key   string `json:"key" jsonschema_description:"Tool key to execute"`
	Input string `json:"input" jsonschema_description:"JSON-encoded arguments for the tool"`
}

func generateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)

	data, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}
	return result
}

var messageSchema = generateSchema[Message]()
