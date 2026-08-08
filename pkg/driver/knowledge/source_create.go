package knowledge

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// IteratorConstructor builds the IKnowledgeIterator for a source driver from
// a config. It backs ISourceDriver.Iterator when set via WithIterator.
type IteratorConstructor func(ctx context.Context, config jsonx.JSONMap) (IKnowledgeIterator, error)

// ReaderConstructor builds the IKnowledgeReader for a single item from a
// config and item ID. It backs ISourceDriver.Reader when set via WithReader.
type ReaderConstructor func(ctx context.Context, config jsonx.JSONMap, itemID string) (IKnowledgeReader, error)

// SourceOption configures a ISourceDriver created via CreateSourceDriver.
type SourceOption func(*sourceDriver)

// WithSourceInfo sets the driver's Info wholesale.
func WithSourceInfo(info driver.Info) SourceOption {
	return func(d *sourceDriver) {
		d.info = info
	}
}

// WithSourceKey sets the driver's registry identity (driver.Info.Key). It's
// required — CreateSourceDriver fails without it — since R.Register derives
// the driver's registry key from it.
func WithSourceKey(key string) SourceOption {
	return func(d *sourceDriver) {
		d.info.Key = key
	}
}

// WithSourceName sets the driver's display name (driver.Info.Name).
func WithSourceName(name string) SourceOption {
	return func(d *sourceDriver) {
		d.info.Name = name
	}
}

// WithSourceDescription sets the driver's human-readable description.
func WithSourceDescription(description string) SourceOption {
	return func(d *sourceDriver) {
		d.info.Description = description
	}
}

// WithSourceIcon sets the driver's display icon.
func WithSourceIcon(icon string) SourceOption {
	return func(d *sourceDriver) {
		d.info.Icon = icon
	}
}

// WithSourceTags sets the driver's classification tags.
func WithSourceTags(tags ...string) SourceOption {
	return func(d *sourceDriver) {
		d.info.Tags = tags
	}
}

// WithSourceConfigurationSchema sets the raw JSON Schema (as a
// jsonx.JSONMap) used to validate config passed to the driver.
func WithSourceConfigurationSchema(schemaJSON []byte) SourceOption {
	return func(d *sourceDriver) {
		schema, err := jsonschema.Load(schemaJSON)
		if err == nil {
			d.info.ConfigurationSchema = schema
		}
	}
}

// WithTestConnection sets the func used to verify a config can connect to
// the source.
func WithTestConnection(fn func(ctx context.Context, config jsonx.JSONMap) error) SourceOption {
	return func(d *sourceDriver) {
		d.testConnection = fn
	}
}

// WithIterator sets the constructor used to build the IKnowledgeIterator
// returned by ISourceDriver.Iterator.
func WithIterator(fn IteratorConstructor) SourceOption {
	return func(d *sourceDriver) {
		d.iterator = fn
	}
}

// WithReader sets the constructor used to build the IKnowledgeReader
// returned by ISourceDriver.Reader.
func WithReader(fn ReaderConstructor) SourceOption {
	return func(d *sourceDriver) {
		d.reader = fn
	}
}

// CreateSourceDriver builds a ISourceDriver from the given SourceOptions.
// Key, TestConnection, Iterator, and Reader are required; CreateSourceDriver
// returns an error instead of a ISourceDriver if any of them is missing, so
// a misconfigured driver never gets registered.
func CreateSourceDriver(opts ...SourceOption) (ISourceDriver, error) {
	d := &sourceDriver{}
	for _, opt := range opts {
		opt(d)
	}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil
}
