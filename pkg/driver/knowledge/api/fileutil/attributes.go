package fileutil

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// BuildAttributes builds kind-specific attributes for a source item.
// Values that require deeper format parsing (language, title, duration,
// etc.) are left unset when they cannot be derived from the raw file alone.
// imageDimensions is only called for SourceItemKindImage and may be nil.
func BuildAttributes(kind knowledge.SourceItemKind, mediaType string, size int64, imageDimensions func() (width, height int)) jsonx.JSONMap {
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
		width, height := 0, 0
		if imageDimensions != nil {
			width, height = imageDimensions()
		}
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
		return jsonx.JSONMap{}
	}

	m, err := jsonx.ToJSONMap(attrs)
	if err != nil {
		return jsonx.JSONMap{}
	}
	return m
}

// ImageDimensions decodes r's image header to obtain pixel dimensions
// without reading the full image data. Supports JPEG, PNG, and GIF.
func ImageDimensions(r io.Reader) (width, height int) {
	cfg, _, err := image.DecodeConfig(r)
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}
