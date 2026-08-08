package openai

import (
	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/api/openai_compatible"
)

const baseURL = "https://api.openai.com/v1"

func New() (llm.Driver, error) {
	return llm.CreateDriver(
		llm.WithKey("openai"),
		llm.WithName("OpenAI"),
		llm.WithDescription("OpenAI language models."),
		llm.WithIcon("https://openai.com/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(openaicompatible.DefaultConfigSchema),
		llm.WithAPI(openaicompatible.New(baseURL)),
	)
}
