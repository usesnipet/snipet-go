package groq

import (
	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/api/openai_compatible"
)

const baseURL = "https://api.groq.com/openai/v1"

func New() (llm.Driver, error) {
	return llm.CreateDriver(
		llm.WithKey("groq"),
		llm.WithName("Groq"),
		llm.WithDescription("Groq high-speed inference models."),
		llm.WithIcon("https://groq.com/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(openaicompatible.DefaultConfigSchema),
		llm.WithAPI(openaicompatible.New(baseURL)),
	)
}
