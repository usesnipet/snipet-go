package llm

import (
	"github.com/usesnipet/snipet/drivers/llm/gemini"
	"github.com/usesnipet/snipet/internal/registry"
	"github.com/usesnipet/snipet/internal/runtime/driver"
)

func Registry() *registry.R[driver.ILLM] {
	registry := registry.New[driver.ILLM]()
	// registry.MustRegister("openai", openai.New())
	registry.MustRegister("gemini", gemini.New())

	return registry
}
