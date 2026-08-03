package llm

import "github.com/usesnipet/snipet/pkg/driver/tool"

type GenerateOptions struct {
	Prompt Prompt

	Tools tool.Toolset
}
