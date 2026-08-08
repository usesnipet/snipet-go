package llm

import (
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// Option configures a Driver created via CreateDriver.
type Option func(*llmDriver)

func WithInfo(info driver.Info) Option {
	return func(o *llmDriver) {
		o.info = info
	}
}

// WithKey sets the driver's registry identity (driver.Info.Key). It's
// required — CreateDriver fails without it — since R.Register derives the
// driver's registry key from it.
func WithKey(key string) Option {
	return func(o *llmDriver) {
		o.info.Key = key
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

// WithConfigurationSchema sets the raw JSON Schema (as a jsonx.JSONMap) used
// to validate config passed to the driver. Prefer ConfigurationSchema or
// MustConfigurationSchema to build this value from a JSON document.
func WithConfigurationSchema(schema jsonx.JSONMap) Option {
	return func(o *llmDriver) {
		o.info.ConfigurationSchema = schema
	}
}

func WithModelLoader(modelLoader ModelLoader) Option {
	return func(o *llmDriver) {
		o.modelLoader = modelLoader
	}
}
