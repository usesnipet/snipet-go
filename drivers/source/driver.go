package source

import (
	"github.com/usesnipet/snipet/drivers/source/fs"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func Registry(log *logger.Logger) *driver.Registry[knowledge.ISourceDriver] {
	r := driver.NewRegistry[knowledge.ISourceDriver](log)
	r.Register(fs.NewDriver(), nil)

	return r
}
