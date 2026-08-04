package openaicompatible

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

// New returns an llm.API that talks to an OpenAI-compatible Chat Completions
// endpoint over plain HTTP. baseURL is the API root (e.g. "https://api.openai.com/v1").
// Config.Endpoint, when set, overrides baseURL at request time.
func New(baseURL string) llm.API {
	return llm.API{
		TestConnection: func(ctx context.Context, config util.JSONMap) error {
			return testConnection(ctx, baseURL, config)
		},
		Stream: func(ctx context.Context, config util.JSONMap, options llm.GenerateOptions) (<-chan llm.StreamEvent, error) {
			return stream(ctx, baseURL, config, options)
		},
	}
}
