package llm

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/jsonx"
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

// Validate checks Info is well-formed and TestConnection, Stream, and
// Generate are all configured. It's called by CreateDriver and again by
// R.Register, so a driver missing any of these never enters a registry.
func (d *llmDriver) Validate() error {
	if err := d.info.Validate(); err != nil {
		return err
	}
	if d.api.TestConnection == nil {
		return ErrTestConnectionNotConfigured
	}
	if d.api.Stream == nil {
		return ErrStreamNotConfigured
	}
	if d.api.Generate == nil {
		return ErrGenerateNotConfigured
	}
	return nil
}

func (d *llmDriver) TestConnection(ctx context.Context, config jsonx.JSONMap) error {
	return d.api.TestConnection(ctx, config)
}

func (d *llmDriver) Stream(ctx context.Context, config jsonx.JSONMap, options GenerateOptions) (StreamIterator, error) {
	return d.api.Stream(ctx, config, options)
}

func (d *llmDriver) Generate(ctx context.Context, config jsonx.JSONMap, options GenerateOptions) (GenerateResult, error) {
	return d.api.Generate(ctx, config, options)
}

func (d *llmDriver) Models(ctx context.Context, config jsonx.JSONMap) ([]Model, error) {
	if d.modelLoader.Models == nil {
		return nil, ErrModelLoaderNotConfigured
	}
	return d.modelLoader.Models(ctx, config)
}

func (d *llmDriver) Model(ctx context.Context, config jsonx.JSONMap) (Model, error) {
	if d.modelLoader.Model == nil {
		return Model{}, ErrModelLoaderNotConfigured
	}
	return d.modelLoader.Model(ctx, config)
}
