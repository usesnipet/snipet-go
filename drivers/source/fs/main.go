package fs

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

//go:embed schema.json
var schemaJSON []byte

type Driver struct{}

func NewDriver() runtime.SourceDriver {
	return &Driver{}
}

type driverConfig struct {
	BasePath string
}

func parseConfig(config util.JSONMap) (driverConfig, error) {
	basePath, ok := config["basePath"].(string)
	if !ok || strings.TrimSpace(basePath) == "" {
		return driverConfig{}, fmt.Errorf("basePath is required")
	}
	return driverConfig{BasePath: basePath}, nil
}

func (d *Driver) GetConfigurationSchema(ctx context.Context) (util.JSONMap, error) {
	var schema util.JSONMap
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return nil, fmt.Errorf("fs: parse schema: %w", err)
	}
	return schema, nil
}

func (d *Driver) TestConnection(ctx context.Context, config util.JSONMap) error {
	cfg, err := parseConfig(config)
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

func (d *Driver) Iterator(ctx context.Context, config util.JSONMap) (runtime.SourceIterator, error) {
	cfg, err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	files, err := listFiles(cfg.BasePath)
	if err != nil {
		return nil, err
	}
	return NewIterator(files), nil
}
