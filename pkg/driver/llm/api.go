package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

// API holds the provider methods used by a driver for connectivity checks and generation.
type API struct {
	TestConnection func(ctx context.Context, config util.JSONMap) error
	Stream         func(ctx context.Context, config util.JSONMap, options GenerateOptions) (<-chan StreamEvent, error)
	Generate       func(ctx context.Context, config util.JSONMap, options GenerateOptions) (GenerateResult, error)
}
