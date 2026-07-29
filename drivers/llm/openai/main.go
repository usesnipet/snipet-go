package openai

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaiprovider "github.com/usesnipet/snipet/pkg/driver/llm/openai"
)

//go:embed schema.json
var schemaJSON []byte

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("OpenAI"),
		llm.WithDescription("OpenAI is a language model provider. Also works with OpenAI-compatible APIs via base_url."),
		llm.WithIcon("https://openai.com/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(openaiprovider.New()),
	)
}
