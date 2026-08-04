package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
)

// llmDriver is the concrete Driver built by CreateDriver, delegating each
// method to the corresponding func in api, if configured.
type llmDriver struct {
	info        driver.Info
	api         API
	modelLoader ModelLoader
}

func (d *llmDriver) Info() driver.Info {
	return d.info
}

func (d *llmDriver) TestConnection(ctx context.Context, config util.JSONMap) error {
	if d.api.TestConnection == nil {
		return ErrTestConnectionNotConfigured
	}
	return d.api.TestConnection(ctx, config)
}

func (d *llmDriver) Stream(ctx context.Context, config util.JSONMap, options GenerateOptions) (<-chan StreamEvent, error) {
	if d.api.Stream == nil {
		return nil, ErrStreamNotConfigured
	}
	return d.api.Stream(ctx, config, options)
}

func (d *llmDriver) Generate(ctx context.Context, config util.JSONMap, options GenerateOptions) (GenerateResult, error) {
	if d.api.Generate == nil {
		return GenerateResult{}, ErrGenerateNotConfigured
	}
	return d.api.Generate(ctx, config, options)
}

func (d *llmDriver) Models(ctx context.Context, config util.JSONMap) ([]Model, error) {
	if d.modelLoader.Models == nil {
		return nil, ErrModelLoaderNotConfigured
	}
	return d.modelLoader.Models(ctx, config)
}

func (d *llmDriver) Model(ctx context.Context, config util.JSONMap) (Model, error) {
	if d.modelLoader.Model == nil {
		return Model{}, ErrModelLoaderNotConfigured
	}
	return d.modelLoader.Model(ctx, config)
}
