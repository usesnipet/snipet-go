package source

import (
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/runtime/driver"
)

func Registry() *registry.R[driver.IKnowledgeSource] {
	registry := registry.New[driver.IKnowledgeSource]()
	// registry.MustRegister("fs", fs.NewDriver())

	return registry
}
