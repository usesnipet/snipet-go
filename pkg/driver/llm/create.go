package llm

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
)

// Option configures a Driver created via CreateDriver.
type Option func(*options)

type options struct {
	info driver.Info
	api  API
}

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

// WithConfigurationSchema sets the raw JSON Schema (as a util.JSONMap) used
// to validate config passed to the driver. Prefer ConfigurationSchema or
// MustConfigurationSchema to build this value from a JSON document.
func WithConfigurationSchema(schema util.JSONMap) Option {
	return func(o *options) {
		o.info.ConfigurationSchema = schema
	}
}

// ConfigurationSchema converts a JSON schema document into a util.JSONMap.
func ConfigurationSchema(schemaJSON []byte) (util.JSONMap, error) {
	return jsonschema.Load(schemaJSON)
}

// MustConfigurationSchema is like ConfigurationSchema but panics on error.
func MustConfigurationSchema(schemaJSON []byte) util.JSONMap {
	schema, err := ConfigurationSchema(schemaJSON)
	if err != nil {
		panic(err)
	}
	return schema
}

// CreateDriver builds a Driver from the given Options. The actual behavior
// (TestConnection/Generate/Stream) comes from WithAPI; any method whose
// underlying func is nil returns an error when called.
func CreateDriver(opts ...Option) Driver {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return &llmDriver{info: o.info, api: o.api}
}

// llmDriver is the concrete Driver built by CreateDriver, delegating each
// method to the corresponding func in api, if configured.
type llmDriver struct {
	info driver.Info
	api  API
}

func (d *llmDriver) Info() driver.Info {
	return d.info
}

func (d *llmDriver) TestConnection(ctx context.Context, config util.JSONMap) error {
	if d.api.TestConnection == nil {
		return fmt.Errorf("test connection not configured")
	}
	return d.api.TestConnection(ctx, config)
}

func (d *llmDriver) Stream(ctx context.Context, config util.JSONMap, options GenerateOptions) (<-chan StreamEvent, error) {
	if d.api.Stream == nil {
		return nil, fmt.Errorf("stream not configured")
	}
	return d.api.Stream(ctx, config, options)
}

func (d *llmDriver) Capabilities(ctx context.Context, config util.JSONMap) (Capabilities, error) {
	if d.api.Capabilities == nil {
		return Capabilities{ToolCall: true}, nil
	}
	return d.api.Capabilities(ctx, config)
}
