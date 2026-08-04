package generator

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

type Strategy interface {
	Generate(ctx context.Context, execution *Execution, toolset tool.Toolset) (message msg.Message, err error)
}
