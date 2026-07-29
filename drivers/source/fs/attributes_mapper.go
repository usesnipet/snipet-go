package fs

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/internal/util"
)

// mapAttributes builds kind-specific attributes for a local file.
// Values that require deeper format parsing (language, title, duration, etc.)
// are left unset when they cannot be derived from the filesystem alone.
func mapAttributes(kind knowledge.SourceItemKind, path string, info os.FileInfo) util.JSONMap {
	mediaType := detectMediaType(path)
	size := info.Size()

	var attrs any
	switch kind {
	case knowledge.SourceItemKindDocument:
		attrs = knowledge.DocumentAttributes{
			MediaType: mediaType,
			Size:      size,
		}
	case knowledge.SourceItemKindText:
		attrs = knowledge.TextAttributes{}
	case knowledge.SourceItemKindImage:
		width, height := imageDimensions(path)
		attrs = knowledge.ImageAttributes{
			MediaType: mediaType,
			Size:      size,
			Width:     width,
			Height:    height,
		}
	case knowledge.SourceItemKindAudio:
		attrs = knowledge.AudioAttributes{
			MediaType: mediaType,
			Size:      size,
		}
	case knowledge.SourceItemKindVideo:
		attrs = knowledge.VideoAttributes{
			MediaType: mediaType,
			Size:      size,
		}
	case knowledge.SourceItemKindStructured:
		attrs = knowledge.StructuredAttributes{
			MediaType: mediaType,
			Size:      size,
		}
	default:
		return util.JSONMap{}
	}

	m, err := util.ToJSONMap(attrs)
	if err != nil {
		return util.JSONMap{}
	}
	return m
}

func imageDimensions(path string) (width, height int) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}
