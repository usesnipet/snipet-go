package index

import (
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/runtime/driver"
)

func Registry() *registry.R[driver.IKnowledgeIndex] {
	registry := registry.New[driver.IKnowledgeIndex]()
	// registry.MustRegister("rag", rag.NewDriver())

	return registry
}
