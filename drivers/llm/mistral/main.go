package mistral

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	gollmprovider "github.com/usesnipet/snipet/pkg/driver/llm/gollm"
)

//go:embed schema.json
var schemaJSON []byte

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("Mistral"),
		llm.WithDescription("Mistral language models via gollm."),
		llm.WithIcon("https://mistral.ai/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(gollmprovider.New("mistral")),
	)
}
