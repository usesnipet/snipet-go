package groq

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/openai_compatible"
)

//go:embed schema.json
var schemaJSON []byte

const baseURL = "https://api.groq.com/openai/v1"

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("Groq"),
		llm.WithDescription("Groq high-speed inference models."),
		llm.WithIcon("https://groq.com/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(openaicompatible.New(baseURL)),
	)
}
