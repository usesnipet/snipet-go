package llm

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/util"
	jsonschema "github.com/usesnipet/snipet/internal/util/json_schema"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/msg"
)

type Option func(*options)

type options struct {
	info driver.Info
	api  API
}

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

func CreateDriver(opts ...Option) Driver {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return &llmDriver{info: o.info, api: o.api}
}

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

func (d *llmDriver) Generate(
	ctx context.Context,
	config util.JSONMap,
	instructions string,
	messages []msg.Message,
) (msg.Message, error) {
	if d.api.Generate == nil {
		return msg.Message{}, fmt.Errorf("generate not configured")
	}
	return d.api.Generate(ctx, config, instructions, messages)
}
