package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFilesReturnsEmptySliceForEmptyDirectory(t *testing.T) {
	base := t.TempDir()

	files, err := listFiles(Config{BasePath: base})

	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestListFilesReturnsAllFilesRecursively(t *testing.T) {
	base := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(base, "root.txt"), []byte("root"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(base, "sub"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(base, "sub", "nested.txt"), []byte("nested"), 0o644))

	files, err := listFiles(Config{BasePath: base})

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{
		filepath.Join(base, "root.txt"),
		filepath.Join(base, "sub", "nested.txt"),
	}, files)
}

func TestListFilesIgnoresMatchingFilesAndDirectories(t *testing.T) {
	base := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(base, "keep.txt"), []byte("keep"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(base, "skip.log"), []byte("skip"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(base, "keep-dir"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(base, "keep-dir", "inside.txt"), []byte("inside"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(base, "ignored-dir"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(base, "ignored-dir", "hidden.txt"), []byte("hidden"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(base, "web", "node_modules", "pkg"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(base, "web", "node_modules", "pkg", "index.js"), []byte("js"), 0o644))

	files, err := listFiles(Config{
		BasePath: base,
		Ignore:   []string{"*.log", "ignored-dir/**", "**/node_modules/**"},
	})

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{
		filepath.Join(base, "keep.txt"),
		filepath.Join(base, "keep-dir", "inside.txt"),
	}, files)
}

func TestListFilesReturnsErrorForMissingBasePath(t *testing.T) {
	files, err := listFiles(Config{BasePath: filepath.Join(t.TempDir(), "missing")})

	require.Error(t, err)
	assert.Nil(t, files)
}
