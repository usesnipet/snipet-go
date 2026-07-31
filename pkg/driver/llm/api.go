package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/msg"
)

// API holds the provider methods used by a driver for connectivity checks and generation.
type API struct {
	TestConnection func(ctx context.Context, config util.JSONMap) error
	Generate       func(ctx context.Context, config util.JSONMap, prompt Prompt) (msg.Message, error)
	Stream         func(ctx context.Context, config util.JSONMap, prompt Prompt) (<-chan StreamDelta, error)
}

func WithAPI(api API) Option {
	return func(o *options) {
		o.api = api
	}
}
