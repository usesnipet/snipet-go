package ollama

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/openai_compatible"
)

//go:embed schema.json
var schemaJSON []byte

const baseURL = "http://localhost:11434/v1"

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("Ollama"),
		llm.WithDescription("Local Ollama models."),
		llm.WithIcon("https://ollama.com/public/icon.png"),
		llm.WithTags("language", "model", "llm", "local"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(openaicompatible.New(baseURL)),
	)
}
