package tool

import (
	"github.com/usesnipet/snipet/drivers/tool/swapi"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver"
	toolDriver "github.com/usesnipet/snipet/pkg/driver/tool"
)

// Registry builds the tool driver registry. A driver that fails to
// construct (e.g. a required option wasn't set) is logged and skipped
// rather than crashing the whole registry.
func Registry(log *logger.Logger) *driver.Registry[toolDriver.Driver] {
	r := driver.NewRegistry[toolDriver.Driver](log)
	r.Register(swapi.New())

	return r
}
