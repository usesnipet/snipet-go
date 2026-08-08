package swapi

import (
	_ "embed"

	"github.com/usesnipet/snipet/pkg/driver/tool"
)

//go:embed tools.json
var toolsJSON []byte

func New() (tool.Driver, error) {
	return tool.CreateDriver(
		tool.WithKey("swapi"),
		tool.WithName("SWAPI"),
		tool.WithDescription("Star Wars API tool"),
		tool.WithIcon("https://swapi.info/favicon.ico"),
		tool.WithToolSetSchema(toolsJSON),
		tool.WithAPI(NewAPI()),
	)
}
