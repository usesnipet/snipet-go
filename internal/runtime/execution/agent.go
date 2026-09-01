package execution

import "github.com/usesnipet/snipet/pkg/jsonx"

type LLMConfig struct {
	Key    string        `json:"key"`
	Config jsonx.JSONMap `json:"config"`
}

type Agent struct {
	Name         string
	Description  string
	Instructions string
	LLM          LLMConfig
}

func NewAgent(name string, description string, instructions string, llm LLMConfig) *Agent {
	return &Agent{
		Name:         name,
		Description:  description,
		Instructions: instructions,
		LLM:          llm,
	}
}
