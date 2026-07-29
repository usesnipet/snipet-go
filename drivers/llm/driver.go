package llm

import (
	"github.com/usesnipet/snipet/drivers/llm/gemini"
	"github.com/usesnipet/snipet/drivers/llm/openai"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	llmDriver "github.com/usesnipet/snipet/pkg/driver/llm"
)

func Registry() *registry.R[llmDriver.Driver] {
	registry := registry.New[llmDriver.Driver]()
	registry.MustRegister("openai", openai.New())
	registry.MustRegister("gemini", gemini.New())

	return registry
}
