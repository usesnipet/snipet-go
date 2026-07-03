package fs_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/drivers/source/fs"
	"github.com/usesnipet/snipet/internal/util"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func hashContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func TestGetConfigurationSchema(t *testing.T) {
	t.Parallel()

	driver := fs.NewDriver()

	schema, err := driver.GetConfigurationSchema(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "object", schema["type"])
	assert.Contains(t, schema["properties"], "basePath")
}

func TestTestConnectionSucceedsForReadableDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	driver := fs.NewDriver()

	err := driver.TestConnection(context.Background(), util.JSONMap{"basePath": dir})
	require.NoError(t, err)
}

func TestTestConnectionFailsForMissingDirectory(t *testing.T) {
	t.Parallel()

	driver := fs.NewDriver()

	err := driver.TestConnection(context.Background(), util.JSONMap{"basePath": filepath.Join(t.TempDir(), "missing")})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "stat base path")
}

func TestTestConnectionFailsWhenBasePathIsFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := writeFile(t, dir, "file.txt", "content")
	driver := fs.NewDriver()

	err := driver.TestConnection(context.Background(), util.JSONMap{"basePath": filePath})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a directory")
}

func TestTestConnectionFailsWithoutBasePath(t *testing.T) {
	t.Parallel()

	driver := fs.NewDriver()

	err := driver.TestConnection(context.Background(), util.JSONMap{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "basePath is required")
}

func TestScanReturnsFilesWithHashAndMetadata(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := "hello world"
	writeFile(t, dir, "docs/readme.txt", content)

	driver := fs.NewDriver()
	items, err := driver.Scan(context.Background(), util.JSONMap{"basePath": dir}, nil, nil)
	require.NoError(t, err)
	require.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "docs/readme.txt", item.ID)
	assert.Equal(t, "readme.txt", item.Name)
	assert.Equal(t, hashContent(content), item.Hash)
	assert.Equal(t, "docs/readme.txt", item.Metadata["path"])
	assert.EqualValues(t, len(content), item.Metadata["size"])
	require.NotNil(t, item.LastModified)
}

func TestScanPaginatesResults(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, dir, "a.txt", "a")
	writeFile(t, dir, "b.txt", "b")
	writeFile(t, dir, "c.txt", "c")

	driver := fs.NewDriver()
	skip := 1
	take := 1

	items, err := driver.Scan(context.Background(), util.JSONMap{"basePath": dir}, &take, &skip)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "b.txt", items[0].Name)
}
