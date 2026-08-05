package llm

import (
	"github.com/usesnipet/snipet/drivers/llm/groq"
	"github.com/usesnipet/snipet/drivers/llm/mistral"
	"github.com/usesnipet/snipet/drivers/llm/ollama"
	"github.com/usesnipet/snipet/drivers/llm/openai"
	"github.com/usesnipet/snipet/drivers/llm/openrouter"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	llmDriver "github.com/usesnipet/snipet/pkg/driver/llm"
)

func Registry() *registry.R[llmDriver.Driver] {
	registry := registry.New[llmDriver.Driver]()
	registry.MustRegister("openai", openai.New())
	registry.MustRegister("groq", groq.New())
	registry.MustRegister("ollama", ollama.New())
	registry.MustRegister("mistral", mistral.New())
	registry.MustRegister("openrouter", openrouter.New())

	return registry
}
