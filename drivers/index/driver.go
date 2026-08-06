package index

import (
	"github.com/usesnipet/snipet/drivers/index/rag"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func Registry() *registry.R[knowledge.IIndexDriver] {
	registry := registry.New[knowledge.IIndexDriver]()
	registry.MustRegister("rag", rag.NewDriver())

	return registry
}
