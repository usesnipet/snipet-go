package tool

import (
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/runtime/driver"
)

func Registry() *registry.R[driver.ITool] {
	registry := registry.New[driver.ITool]()

	return registry
}
