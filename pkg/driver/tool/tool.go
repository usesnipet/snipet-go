package tool

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
)

type Driver interface {
	driver.IDriver

	Execute(ctx context.Context, config util.JSONMap, call Call) Result
}

type Call struct {
	Key   string
	Input any
}

type Result struct {
	Key     string
	Success bool
	Output  any
	Error   error
}
