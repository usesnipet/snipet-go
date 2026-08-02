package tool

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
)

type Driver interface {
	driver.IDriver

	ToolSet() Toolset
	Call(ctx context.Context, call ToolCall) (ToolResult, error)
}

type ToolCall struct {
	Tool      string         `json:"tool"`
	Arguments map[string]any `json:"arguments"`
}

type ToolResult struct {
	Tool      string         `json:"tool"`
	Arguments map[string]any `json:"arguments"`
	Result    string         `json:"result"`
}
