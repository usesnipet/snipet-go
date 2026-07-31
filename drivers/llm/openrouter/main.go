package openrouter

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	gollmprovider "github.com/usesnipet/snipet/pkg/driver/llm/gollm"
)

//go:embed schema.json
var schemaJSON []byte

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("OpenRouter"),
		llm.WithDescription("OpenRouter multi-provider models via gollm."),
		llm.WithIcon("https://openrouter.ai/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(gollmprovider.New("openrouter")),
	)
}
