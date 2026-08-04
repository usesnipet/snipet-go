package static

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

func New(models []llm.Model) llm.ModelLoader {
	return llm.ModelLoader{
		Models: func(ctx context.Context, config util.JSONMap) ([]llm.Model, error) {
			return models, nil
		},
		Model: func(ctx context.Context, config util.JSONMap, name string) (llm.Model, error) {
			for _, model := range models {
				if model.Name == name {
					return model, nil
				}
			}
			return llm.Model{}, llm.ErrModelNotFound
		},
	}
}
