package tool

import (
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

func Registry() *registry.R[tool.Driver] {
	registry := registry.New[tool.Driver]()

	return registry
}
