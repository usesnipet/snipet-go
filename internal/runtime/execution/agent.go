package execution

import "github.com/usesnipet/snipet/internal/runtime/manager"

type LLMConfig manager.Configuration

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
