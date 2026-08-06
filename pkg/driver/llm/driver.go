// Package llm defines the driver contract for LLM provider integrations
// (Driver), the inputs/outputs of a generation call (GenerateOptions,
// GenerateResult, Prompt), streaming (StreamIterator, StreamEvent) and model
// discovery (Model, ModelLoader). Concrete providers are built with
// CreateDriver from an API implementation (see the api subpackages, e.g.
// openai_compatible).
package llm

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// GenerateOptions carries the per-call inputs to Driver.Generate and
// Driver.Stream: the Prompt to send and the Tools the model may call.
type GenerateOptions struct {
	Prompt Prompt

	Tools tool.Toolset
}

// GenerateResult is the output of Driver.Generate.
type GenerateResult struct {
	Text      string
	ToolCalls []tool.Call
}

// Driver is the contract implemented by every LLM provider integration.
// Generate performs a single blocking completion; Stream performs the same
// call but delivers incremental StreamEvent values as they arrive.
type Driver interface {
	driver.IDriver

	Stream(ctx context.Context, config jsonx.JSONMap, options GenerateOptions) (StreamIterator, error)
	Generate(ctx context.Context, config jsonx.JSONMap, options GenerateOptions) (GenerateResult, error)

	Models(ctx context.Context, config jsonx.JSONMap) ([]Model, error)
	Model(ctx context.Context, config jsonx.JSONMap) (Model, error)
}
