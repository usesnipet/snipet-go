package llm

import (
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
)

// Option configures a Driver created via CreateDriver.
type Option func(*llmDriver)

func WithInfo(info driver.Info) Option {
	return func(o *llmDriver) {
		o.info = info
	}
}

// WithName sets the driver's display name (driver.Info.Name).
func WithName(name string) Option {
	return func(o *llmDriver) {
		o.info.Name = name
	}
}

// WithAPI sets the driver's API.
func WithAPI(api API) Option {
	return func(o *llmDriver) {
		o.api = api
	}
}

// WithDescription sets the driver's human-readable description.
func WithDescription(description string) Option {
	return func(o *llmDriver) {
		o.info.Description = description
	}
}

// WithIcon sets the driver's display icon.
func WithIcon(icon string) Option {
	return func(o *llmDriver) {
		o.info.Icon = icon
	}
}

// WithTags sets the driver's classification tags.
func WithTags(tags ...string) Option {
	return func(o *llmDriver) {
		o.info.Tags = tags
	}
}

// WithConfigurationSchema sets the raw JSON Schema (as a util.JSONMap) used
// to validate config passed to the driver. Prefer ConfigurationSchema or
// MustConfigurationSchema to build this value from a JSON document.
func WithConfigurationSchema(schema util.JSONMap) Option {
	return func(o *llmDriver) {
		o.info.ConfigurationSchema = schema
	}
}

func WithModelLoader(modelLoader ModelLoader) Option {
	return func(o *llmDriver) {
		o.modelLoader = modelLoader
	}
}
