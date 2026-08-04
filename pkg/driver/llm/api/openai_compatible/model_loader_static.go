package openaicompatible

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

func NewStaticModelLoader(models []llm.Model) llm.ModelLoader {
	return llm.ModelLoader{
		Models: func(ctx context.Context, config util.JSONMap) ([]llm.Model, error) {
			return models, nil
		},
		Model: func(ctx context.Context, config util.JSONMap) (llm.Model, error) {
			cfg, err := NewConfig(config)
			if err != nil {
				return llm.Model{}, err
			}

			for _, model := range models {
				if model.Name == cfg.Model {
					return model, nil
				}
			}
			return llm.Model{}, llm.ErrModelNotFound
		},
	}
}
