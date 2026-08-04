package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
)

// Driver is the contract implemented by every LLM provider integration.
// Generate performs a single blocking completion; Stream performs the same
// call but delivers incremental StreamEvent values as they arrive.
type Driver interface {
	driver.IDriver

	Stream(ctx context.Context, config util.JSONMap, options GenerateOptions) (<-chan StreamEvent, error)

	// Capabilities reports what the configured model supports (see
	// Capabilities). Drivers that don't need per-config detection (i.e.
	// don't set API.Capabilities) report ToolCall: true.
	Capabilities(ctx context.Context, config util.JSONMap) (Capabilities, error)
}
