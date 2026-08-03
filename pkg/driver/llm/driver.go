package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/msg"
)

type Driver interface {
	driver.IDriver

	Generate(ctx context.Context, config util.JSONMap, options GenerateOptions) (msg.Message, error)
	Stream(ctx context.Context, config util.JSONMap, options GenerateOptions) (<-chan StreamEvent, error)
}
