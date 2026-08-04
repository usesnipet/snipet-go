package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

type ModelLoader struct {
	Models func(ctx context.Context, config util.JSONMap) ([]Model, error)
	Model  func(ctx context.Context, config util.JSONMap) (Model, error)
}
