package tool

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
)

type options struct {
	info    driver.Info
	toolset Toolset
	api     API
}

type Option func(*options)

func WithName(name string) Option {
	return func(o *options) {
		o.info.Name = name
	}
}

func WithDescription(description string) Option {
	return func(o *options) {
		o.info.Description = description
	}
}

func WithIcon(icon string) Option {
	return func(o *options) {
		o.info.Icon = icon
	}
}

func WithTags(tags ...string) Option {
	return func(o *options) {
		o.info.Tags = tags
	}
}

func WithConfigurationSchema(schema util.JSONMap) Option {
	return func(o *options) {
		o.info.ConfigurationSchema = schema
	}
}

func WithToolSet(toolset Toolset) Option {
	return func(o *options) {
		o.toolset = toolset
	}
}

func WithToolSetSchema(schemaJSON []byte) Option {
	return func(o *options) {
		toolSet, err := ToolSetFromSchema(schemaJSON)
		if err != nil {
			panic(err)
		}
		o.toolset = toolSet
	}
}

func CreateDriver(opts ...Option) Driver {
	o := &options{
		toolset: Toolset{},
		info:    driver.Info{},
	}
	for _, opt := range opts {
		opt(o)
	}

	return &toolDriver{info: o.info, toolset: o.toolset, api: o.api}
}

type toolDriver struct {
	info    driver.Info
	toolset Toolset
	api     API
}

func (d *toolDriver) Info() driver.Info {
	return d.info
}

func (d *toolDriver) TestConnection(ctx context.Context, config util.JSONMap) error {
	if d.api.TestConnection == nil {
		return fmt.Errorf("test connection not configured")
	}
	return d.api.TestConnection(ctx, config)
}

func (d *toolDriver) ToolSet() Toolset {
	return d.toolset
}

func (d *toolDriver) Call(ctx context.Context, call ToolCall) (ToolResult, error) {
	if d.api.Call == nil {
		return ToolResult{}, fmt.Errorf("call not configured")
	}
	return d.api.Call(ctx, call)
}
