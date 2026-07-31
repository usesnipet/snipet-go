package gollm

// structuredMessage is the structured output the model must return.
type structuredMessage struct {
	Content string `json:"content"`
}

const outputSpec = `Respond with a JSON object matching this schema: {"type":"object","properties":{"content":{"type":"string","description":"Assistant text content for this turn"}},"required":["content"],"additionalProperties":false}`
