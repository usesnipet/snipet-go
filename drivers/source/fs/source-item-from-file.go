package fs

import (
	"os"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/driver/knowledge/api/fileutil"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// sourceItemFromFile builds a SourceItem from a local file path.
// The item ID is the absolute path, Kind and Attributes are inferred from
// content/extension, and Metadata includes size, last_modified, path, and
// name. The second return value is a content hash used for change detection.
func sourceItemFromFile(path string) (*knowledge.SourceItem, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	kind, attributes, rest := fileutil.Inspect(path, info.Size(), file)

	hash, err := fileutil.Hash(rest)
	if err != nil {
		return nil, "", err
	}

	lastModified := info.ModTime()
	item := &knowledge.SourceItem{
		ID:           path,
		Name:         info.Name(),
		Kind:         kind,
		LastModified: &lastModified,
		Attributes:   attributes,
		Metadata: jsonx.JSONMap{
			"size":          info.Size(),
			"last_modified": lastModified,
			"path":          path,
			"name":          info.Name(),
		},
	}
	return item, hash, nil
}
