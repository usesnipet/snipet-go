package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

// Capabilities describes what a driver supports for a given config, e.g.
// the specific model selected in it. Providers whose catalog isn't uniform
// (self-hosted ones like Ollama, where installed models vary in what they
// support) report it per config instead of assuming a fixed answer.
type Capabilities struct {
	// ToolCall reports whether the configured model accepts the Tools field
	// of GenerateOptions natively.
	ToolCall bool
}

// API holds the provider methods used by a driver for connectivity checks and generation.
type API struct {
	TestConnection func(ctx context.Context, config util.JSONMap) error
	Stream         func(ctx context.Context, config util.JSONMap, options GenerateOptions) (<-chan StreamEvent, error)

	// Capabilities reports what config supports. Optional: when nil, the
	// driver is assumed to fully support tool calls, which holds for every
	// hosted provider's catalog today.
	Capabilities func(ctx context.Context, config util.JSONMap) (Capabilities, error)
}

func WithAPI(api API) Option {
	return func(o *options) {
		o.api = api
	}
}
