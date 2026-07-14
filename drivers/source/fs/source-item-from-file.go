package fs

import (
	"os"

	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

// sourceItemFromFile builds a SourceItem from a local file path.
// The item ID is the absolute path, Kind is inferred from content/extension,
// and Metadata includes size, last_modified, path, and name.
// The second return value is a content hash used for change detection.
func sourceItemFromFile(path string) (*runtime.SourceItem, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", err
	}

	hash, err := fileHash(path)
	if err != nil {
		return nil, "", err
	}

	kind := mapKind(path)

	lastModified := info.ModTime()
	item := &runtime.SourceItem{
		ID:           path,
		Name:         info.Name(),
		Kind:         kind,
		LastModified: &lastModified,
		Attributes:   mapAttributes(kind, path, info),
		Metadata: util.JSONMap{
			"size":          info.Size(),
			"last_modified": lastModified,
			"path":          path,
			"name":          info.Name(),
		},
	}
	return item, hash, nil
}
