package tool

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type options struct {
	info    driver.Info
	toolset Toolset
	api     API
}

// Option configures a Driver created via CreateDriver.
type Option func(*options)

// WithName sets the driver's display name (driver.Info.Name).
func WithName(name string) Option {
	return func(o *options) {
		o.info.Name = name
	}
}

// WithDescription sets the driver's human-readable description.
func WithDescription(description string) Option {
	return func(o *options) {
		o.info.Description = description
	}
}

// WithIcon sets the driver's display icon.
func WithIcon(icon string) Option {
	return func(o *options) {
		o.info.Icon = icon
	}
}

// WithTags sets the driver's classification tags.
func WithTags(tags ...string) Option {
	return func(o *options) {
		o.info.Tags = tags
	}
}

// WithConfigurationSchema sets the raw JSON Schema (as a jsonx.JSONMap) used
// to validate config passed to the driver.
func WithConfigurationSchema(schema jsonx.JSONMap) Option {
	return func(o *options) {
		o.info.ConfigurationSchema = schema
	}
}

// WithToolSet sets the Toolset the driver exposes directly.
func WithToolSet(toolset Toolset) Option {
	return func(o *options) {
		o.toolset = toolset
	}
}

// WithToolSetSchema sets the driver's Toolset by parsing and validating
// schemaJSON against the embedded tool-set schema (see ToolSetFromSchema).
// It panics if schemaJSON is invalid, since toolsets are normally built
// from schemas known at compile time.
func WithToolSetSchema(schemaJSON []byte) Option {
	return func(o *options) {
		toolSet, err := ToolSetFromSchema(schemaJSON)
		if err != nil {
			panic(err)
		}
		o.toolset = toolSet
	}
}

// CreateDriver builds a Driver from the given Options. The actual behavior
// (TestConnection/Call) comes from WithAPI; any method whose underlying func
// is nil returns an error when called.
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

// toolDriver is the concrete Driver built by CreateDriver, delegating each
// method to the corresponding func in api, if configured.
type toolDriver struct {
	info    driver.Info
	toolset Toolset
	api     API
}

func (d *toolDriver) Info() driver.Info {
	return d.info
}

func (d *toolDriver) TestConnection(ctx context.Context, config jsonx.JSONMap) error {
	if d.api.TestConnection == nil {
		return fmt.Errorf("test connection not configured")
	}
	return d.api.TestConnection(ctx, config)
}

func (d *toolDriver) ToolSet() Toolset {
	return d.toolset
}

func (d *toolDriver) Call(ctx context.Context, call Call) (Result, error) {
	if d.api.Call == nil {
		return Result{}, fmt.Errorf("call not configured")
	}
	return d.api.Call(ctx, call)
}
