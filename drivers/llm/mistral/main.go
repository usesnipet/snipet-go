package mistral

import (
	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/api/openai_compatible"
)

const baseURL = "https://api.mistral.ai/v1"

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("Mistral"),
		llm.WithDescription("Mistral language models."),
		llm.WithIcon("https://mistral.ai/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(openaicompatible.DefaultConfigSchema),
		llm.WithAPI(openaicompatible.New(baseURL)),
	)
}
