package fs

import (
	_ "embed"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

func (d *Driver) Scan(ctx context.Context, config util.JSONMap, take *int, skip *int) ([]runtime.SourceItem, error) {
	cfg, err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	items, err := collectItems(ctx, cfg.BasePath)
	if err != nil {
		return nil, err
	}

	return paginateItems(items, skip, take), nil
}

func collectItems(ctx context.Context, basePath string) ([]runtime.SourceItem, error) {
	var items []runtime.SourceItem

	err := filepath.WalkDir(basePath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("fs: stat %q: %w", path, err)
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return fmt.Errorf("fs: relative path for %q: %w", path, err)
		}

		hash, err := hashFile(path)
		if err != nil {
			return fmt.Errorf("fs: hash %q: %w", path, err)
		}

		modTime := info.ModTime()
		items = append(items, runtime.SourceItem{
			ID:   filepath.ToSlash(relPath),
			Name: filepath.Base(path),
			Hash: hash,
			Metadata: util.JSONMap{
				"path": filepath.ToSlash(relPath),
				"size": info.Size(),
			},
			LastModified: &modTime,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items, nil
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func paginateItems(items []runtime.SourceItem, skip, take *int) []runtime.SourceItem {
	start := 0
	if skip != nil {
		start = *skip
	}
	if start >= len(items) {
		return []runtime.SourceItem{}
	}

	items = items[start:]
	if take != nil && *take < len(items) {
		items = items[:*take]
	}

	return items
}
