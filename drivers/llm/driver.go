package llm

import (
	"github.com/usesnipet/snipet/drivers/llm/groq"
	"github.com/usesnipet/snipet/drivers/llm/mistral"
	"github.com/usesnipet/snipet/drivers/llm/ollama"
	"github.com/usesnipet/snipet/drivers/llm/openai"
	"github.com/usesnipet/snipet/drivers/llm/openrouter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver"
	llmDriver "github.com/usesnipet/snipet/pkg/driver/llm"
)

// Registry builds the LLM driver registry. A driver that fails to
// construct (e.g. a required option wasn't set) is logged and skipped
// rather than crashing the whole registry.
func Registry(log *logger.Logger) *driver.Registry[llmDriver.Driver] {
	r := driver.NewRegistry[llmDriver.Driver](log)

	r.Register(openai.New())
	r.Register(groq.New())
	r.Register(ollama.New())
	r.Register(mistral.New())
	r.Register(openrouter.New())

	return r
}
