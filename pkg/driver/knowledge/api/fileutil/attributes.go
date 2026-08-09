package fileutil

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// NewDocumentAttributes builds attributes for a SourceItemKindDocument item.
// language, title, and author are best-effort bibliographic metadata (e.g.
// parsed from PDF/DOCX properties); pass nil when unavailable.
func NewDocumentAttributes(mediaType string, size int64, language, title, author *string) knowledge.DocumentAttributes {
	return knowledge.DocumentAttributes{
		MediaType: mediaType,
		Size:      size,
		Language:  language,
		Title:     title,
		Author:    author,
	}
}

// NewTextAttributes builds attributes for a SourceItemKindText item.
// language, title, and author are best-effort bibliographic metadata; pass
// nil when unavailable.
func NewTextAttributes(mediaType string, size int64, language, title, author *string) knowledge.TextAttributes {
	return knowledge.TextAttributes{
		MediaType: mediaType,
		Size:      size,
		Language:  language,
		Title:     title,
		Author:    author,
	}
}

// NewImageAttributes builds attributes for a SourceItemKindImage item.
func NewImageAttributes(mediaType string, size int64, width, height int) knowledge.ImageAttributes {
	return knowledge.ImageAttributes{
		MediaType: mediaType,
		Size:      size,
		Width:     width,
		Height:    height,
	}
}

// NewAudioAttributes builds attributes for a SourceItemKindAudio item.
// duration (seconds), bitrate, sampleRate, and channels require probing the
// audio stream itself; pass nil when unavailable.
func NewAudioAttributes(mediaType string, size int64, duration, bitrate, sampleRate, channels *int) knowledge.AudioAttributes {
	return knowledge.AudioAttributes{
		MediaType:  mediaType,
		Size:       size,
		Duration:   duration,
		Bitrate:    bitrate,
		SampleRate: sampleRate,
		Channels:   channels,
	}
}

// NewVideoAttributes builds attributes for a SourceItemKindVideo item.
// duration (seconds), bitrate, frameRate, width, and height require probing
// the video stream itself; pass nil when unavailable.
func NewVideoAttributes(mediaType string, size int64, duration, bitrate, frameRate *int, width, height *int) knowledge.VideoAttributes {
	return knowledge.VideoAttributes{
		MediaType: mediaType,
		Size:      size,
		Duration:  duration,
		Bitrate:   bitrate,
		FrameRate: frameRate,
		Width:     width,
		Height:    height,
	}
}

// NewStructuredAttributes builds attributes for a SourceItemKindStructured item.
func NewStructuredAttributes(mediaType string, size int64) knowledge.StructuredAttributes {
	return knowledge.StructuredAttributes{
		MediaType: mediaType,
		Size:      size,
	}
}

// ToJSONMap converts a typed attributes struct (e.g. one of the New*Attributes
// results above) into the generic map stored on a SourceItem, returning an
// empty map if the conversion fails.
func ToJSONMap(attrs any) jsonx.JSONMap {
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

// Inspect sniffs r's leading bytes once to classify the item's kind and
// build its matching attributes. It hands back a reader that replays those
// consumed bytes ahead of the rest of r, so a caller can still hash or
// otherwise stream the full content afterward without a second read (and,
// for a remote-backed source, without a second network round trip).
//
// Fields that require deeper format parsing (bibliographic metadata for
// Document/Text, stream probing for Audio/Video) are left nil — Inspect
// only has a raw byte stream to work with, not a format-specific parser.
func Inspect(name string, size int64, r io.Reader) (knowledge.SourceItemKind, jsonx.JSONMap, io.Reader) {
	header, err := readHeader(r)
	rest := io.MultiReader(bytes.NewReader(header), r)
	if err != nil {
		return knowledge.SourceItemKindUnknown, jsonx.JSONMap{}, rest
	}

	mediaType := mediaTypeFromHeader(name, header)
	kind := MapKind(mediaType)

	var attrs any
	switch kind {
	case knowledge.SourceItemKindDocument:
		attrs = NewDocumentAttributes(mediaType, size, nil, nil, nil)
	case knowledge.SourceItemKindText:
		attrs = NewTextAttributes(mediaType, size, nil, nil, nil)
	case knowledge.SourceItemKindImage:
		width, height := ImageDimensions(bytes.NewReader(header))
		attrs = NewImageAttributes(mediaType, size, width, height)
	case knowledge.SourceItemKindAudio:
		attrs = NewAudioAttributes(mediaType, size, nil, nil, nil, nil)
	case knowledge.SourceItemKindVideo:
		attrs = NewVideoAttributes(mediaType, size, nil, nil, nil, nil, nil)
	case knowledge.SourceItemKindStructured:
		attrs = NewStructuredAttributes(mediaType, size)
	default:
		return knowledge.SourceItemKindUnknown, jsonx.JSONMap{}, rest
	}

	return kind, ToJSONMap(attrs), rest
}
