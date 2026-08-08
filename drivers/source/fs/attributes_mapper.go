package fs

import (
	"os"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/driver/knowledge/api/fileutil"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// mapAttributes builds kind-specific attributes for a local file.
func mapAttributes(kind knowledge.SourceItemKind, mediaType string, path string, info os.FileInfo) jsonx.JSONMap {
	return fileutil.BuildAttributes(kind, mediaType, info.Size(), func() (width, height int) {
		return fileImageDimensions(path)
	})
}

func fileImageDimensions(path string) (width, height int) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0
	}
	defer file.Close()

	return fileutil.ImageDimensions(file)
}
