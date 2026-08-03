package llm

import "github.com/usesnipet/snipet/pkg/driver/tool"

// GenerateOptions carries the per-call inputs to Driver.Generate and
// Driver.Stream: the Prompt to send and the Tools the model may call.
type GenerateOptions struct {
	Prompt Prompt

	Tools tool.Toolset
}
