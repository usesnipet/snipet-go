package openai

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/openai_compatible"
)

//go:embed schema.json
var schemaJSON []byte

const baseURL = "https://api.openai.com/v1"

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("OpenAI"),
		llm.WithDescription("OpenAI language models."),
		llm.WithIcon("https://openai.com/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(openaicompatible.New(baseURL)),
	)
}
