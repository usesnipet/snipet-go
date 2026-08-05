package source

import (
	"github.com/usesnipet/snipet/drivers/source/fs"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func Registry() *registry.R[knowledge.ISourceDriver] {
	registry := registry.New[knowledge.ISourceDriver]()
	registry.MustRegister("fs", fs.NewDriver())

	return registry
}
