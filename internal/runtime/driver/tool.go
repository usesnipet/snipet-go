package driver

import (
	"context"

	"github.com/usesnipet/snipet/internal/runtime/tool"
	"github.com/usesnipet/snipet/internal/util"
)

type ITool interface {
	IDriver
	Execute(ctx context.Context, config util.JSONMap, call tool.Call) tool.Result
}
