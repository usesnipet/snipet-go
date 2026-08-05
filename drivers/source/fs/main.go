package fs

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

//go:embed schema.json
var schemaJSON []byte

type Driver struct{}

func NewDriver() knowledge.ISourceDriver {
	return &Driver{}
}

func (d *Driver) Info() driver.Info {
	schema, err := d.configurationSchema()
	if err != nil {
		panic(err)
	}
	return driver.Info{
		Name:                "Filesystem",
		Description:         "Reads knowledge items from a local filesystem path.",
		Tags:                []string{"source", "filesystem"},
		ConfigurationSchema: schema,
	}
}

func (d *Driver) configurationSchema() (jsonx.JSONMap, error) {
	var schema jsonx.JSONMap
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return nil, fmt.Errorf("fs: parse schema: %w", err)
	}
	return schema, nil
}

func (d *Driver) TestConnection(ctx context.Context, config jsonx.JSONMap) error {
	cfg, err := jsonx.ParseJSONMap[Config](config)
	if err != nil {
		return err
	}

	info, err := os.Stat(cfg.BasePath)
	if err != nil {
		return fmt.Errorf("fs: stat base path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("fs: base path is not a directory")
	}

	if _, err := os.ReadDir(cfg.BasePath); err != nil {
		return fmt.Errorf("fs: read base path: %w", err)
	}

	return nil
}

func (d *Driver) Iterator(ctx context.Context, config jsonx.JSONMap) (knowledge.IKnowledgeIterator, error) {
	cfg, err := jsonx.ParseJSONMap[Config](config)
	if err != nil {
		return nil, err
	}

	files, err := listFiles(cfg)
	if err != nil {
		return nil, err
	}
	return NewIterator(files), nil
}

func (d *Driver) Reader(ctx context.Context, config jsonx.JSONMap, itemID string) (knowledge.IKnowledgeReader, error) {
	return NewReader(config, itemID)
}
