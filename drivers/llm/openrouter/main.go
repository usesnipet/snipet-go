package openrouter

import (
	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/api/openai_compatible"
)

const baseURL = "https://openrouter.ai/api/v1"

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("OpenRouter"),
		llm.WithDescription("OpenRouter multi-provider models."),
		llm.WithIcon("https://openrouter.ai/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(openaicompatible.DefaultConfigSchema),
		llm.WithAPI(openaicompatible.New(baseURL)),
	)
}
