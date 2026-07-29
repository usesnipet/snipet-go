package gemini

import (
	_ "embed"

	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	geminiprovider "github.com/usesnipet/snipet/pkg/driver/llm/gemini"
)

//go:embed schema.json
var schemaJSON []byte

func New() llm.Driver {
	return llm.CreateDriver(
		llm.WithName("Gemini"),
		llm.WithDescription("Gemini is a language model that can generate text, images, and audio."),
		llm.WithIcon("https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"),
		llm.WithTags("language", "model", "llm"),
		llm.WithConfigurationSchema(llm.MustConfigurationSchema(schemaJSON)),
		llm.WithAPI(geminiprovider.New()),
	)
}
