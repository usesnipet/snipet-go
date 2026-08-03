package tool

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
)

type Driver interface {
	driver.IDriver

	ToolSet() Toolset
	Call(ctx context.Context, call Call) (Result, error)
}

type Call struct {
	ID        string         `json:"id"`
	Tool      string         `json:"tool"`
	Arguments map[string]any `json:"arguments"`
}

type Result struct {
	Tool      string         `json:"tool"`
	Arguments map[string]any `json:"arguments"`
	Result    string         `json:"result"`
}
