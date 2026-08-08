package index

import (
	"github.com/usesnipet/snipet/drivers/index/rag"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func Registry(log *logger.Logger) *driver.Registry[knowledge.IIndexDriver] {
	r := driver.NewRegistry[knowledge.IIndexDriver](log)
	r.Register(rag.NewDriver(), nil)

	return r
}
