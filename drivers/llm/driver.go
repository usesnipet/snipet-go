package llm

import (
	"github.com/usesnipet/snipet/drivers/llm/gemini"
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/runtime/registry"
)

func Registry() *registry.R[driver.ILLM] {
	registry := registry.New[driver.ILLM]()
	// registry.MustRegister("openai", openai.New())
	registry.MustRegister("gemini", gemini.New())

	return registry
}
