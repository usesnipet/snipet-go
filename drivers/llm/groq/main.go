package groq

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	gollmprovider "github.com/usesnipet/snipet/pkg/driver/llm/gollm"
)

//go:embed schema.json
var schemaJSON []byte

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("Groq"),
		llm.WithDescription("Groq high-speed inference models via gollm."),
		llm.WithIcon("https://groq.com/favicon.ico"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(gollmprovider.New("groq")),
	)
}
