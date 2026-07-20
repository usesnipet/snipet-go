package runtime

import (
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/util"
)

type LLMConfig driver.Configuration

type ToolConfig map[string]util.JSONMap

type Agent struct {
	Name         string
	Description  string
	Instructions string
	Tools        ToolConfig
	LLMs         []LLMConfig
}

func NewAgent(name string, description string, instructions string, tools ToolConfig, llms []LLMConfig) Agent {
	return Agent{
		Name:         name,
		Description:  description,
		Instructions: instructions,
		Tools:        tools,
		LLMs:         llms,
	}
}
