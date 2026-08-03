package tool

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
)

// Driver is the contract implemented by every tool provider integration.
// ToolSet reports the tools it exposes; Call invokes one of them.
type Driver interface {
	driver.IDriver

	ToolSet() Toolset
	Call(ctx context.Context, call Call) (Result, error)
}

// Call is a request to invoke one tool by name (Tool) with the given
// Arguments, identified by ID so its result (Result) and any related
// msg.Message can be correlated back to it.
type Call struct {
	ID        string         `json:"id"`
	Tool      string         `json:"tool"`
	Arguments map[string]any `json:"arguments"`
}

// Result is the outcome of a Call: the tool invoked, the arguments it ran
// with, and its output.
type Result struct {
	Tool      string         `json:"tool"`
	Arguments map[string]any `json:"arguments"`
	Result    string         `json:"result"`
}
