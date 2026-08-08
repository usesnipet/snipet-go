package fs

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

//go:embed schema.json
var schemaJSON []byte

func NewDriver() (knowledge.ISourceDriver, error) {
	return knowledge.CreateSourceDriver(
		knowledge.WithSourceKey("fs"),
		knowledge.WithSourceName("Filesystem"),
		knowledge.WithSourceDescription("Reads knowledge items from a local filesystem path."),
		knowledge.WithSourceTags("source", "filesystem"),
		knowledge.WithSourceConfigurationSchema(schemaJSON),
		knowledge.WithTestConnection(testConnection),
		knowledge.WithIterator(iterator),
		knowledge.WithReader(reader),
	)
}

func testConnection(ctx context.Context, config jsonx.JSONMap) error {
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

func iterator(ctx context.Context, config jsonx.JSONMap) (knowledge.IKnowledgeIterator, error) {
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

func reader(ctx context.Context, config jsonx.JSONMap, itemID string) (knowledge.IKnowledgeReader, error) {
	return NewReader(config, itemID)
}
