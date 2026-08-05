package ollama

import (
	_ "embed"

	"github.com/usesnipet/snipet/pkg/driver"
	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/api/openai_compatible"
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
		llm.WithConfigurationSchema(driver.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(openaicompatible.New(baseURL)),
		llm.WithModelLoader(modelLoader),
	)
}
