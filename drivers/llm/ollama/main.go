package ollama

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	gollmprovider "github.com/usesnipet/snipet/pkg/driver/llm/gollm"
)

//go:embed schema.json
var schemaJSON []byte

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("Ollama"),
		llm.WithDescription("Local Ollama models via gollm."),
		llm.WithIcon("https://ollama.com/public/icon.png"),
		llm.WithTags("language", "model", "llm", "local"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(gollmprovider.New("ollama")),
	)
}
