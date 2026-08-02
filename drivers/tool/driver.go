package tool

import (
	"github.com/usesnipet/snipet/drivers/tool/swapi"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	toolDriver "github.com/usesnipet/snipet/pkg/driver/tool"
)

func Registry() *registry.R[toolDriver.Driver] {
	registry := registry.New[toolDriver.Driver]()
	registry.MustRegister("swapi", swapi.New())

	return registry
}
