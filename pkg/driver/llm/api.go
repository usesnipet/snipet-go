package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/util"
)

// API holds the provider methods used by a driver for connectivity checks and generation.
type API struct {
	TestConnection func(ctx context.Context, config util.JSONMap) error
	Generate       func(ctx context.Context, config util.JSONMap, instructions string, messages []message.Message) (message.Message, error)
}

func WithAPI(api API) Option {
	return func(o *options) {
		o.api = api
	}
}
