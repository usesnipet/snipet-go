package fileutil

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func TestBuildAttributesDocument(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindDocument, "application/pdf", 1024, nil)

	assert.Equal(t, "application/pdf", attrs["media_type"])
	assert.Equal(t, float64(1024), attrs["size"])
}

func TestBuildAttributesText(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindText, "text/plain", 10, nil)

	assert.NotNil(t, attrs)
	_, hasMediaType := attrs["media_type"]
	assert.False(t, hasMediaType)
}

func TestBuildAttributesImageCallsDimensionsCallback(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindImage, "image/png", 2048, func() (int, int) {
		return 100, 200
	})

	assert.Equal(t, "image/png", attrs["media_type"])
	assert.Equal(t, float64(2048), attrs["size"])
	assert.Equal(t, float64(100), attrs["width"])
	assert.Equal(t, float64(200), attrs["height"])
}

func TestBuildAttributesImageNilCallbackDefaultsToZero(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindImage, "image/png", 2048, nil)

	assert.Equal(t, float64(0), attrs["width"])
	assert.Equal(t, float64(0), attrs["height"])
}

func TestBuildAttributesAudio(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindAudio, "audio/mpeg", 512, nil)

	assert.Equal(t, "audio/mpeg", attrs["media_type"])
	assert.Equal(t, float64(512), attrs["size"])
}

func TestBuildAttributesVideo(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindVideo, "video/mp4", 4096, nil)

	assert.Equal(t, "video/mp4", attrs["media_type"])
}

func TestBuildAttributesStructured(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindStructured, "application/json", 64, nil)

	assert.Equal(t, "application/json", attrs["media_type"])
}

func TestBuildAttributesUnknownKindReturnsEmptyMap(t *testing.T) {
	attrs := BuildAttributes(knowledge.SourceItemKindUnknown, "", 0, nil)

	assert.Empty(t, attrs)
}

func TestImageDimensionsDecodesPNG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 32))
	for x := range 64 {
		for y := range 32 {
			img.Set(x, y, color.White)
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))

	width, height := ImageDimensions(&buf)

	assert.Equal(t, 64, width)
	assert.Equal(t, 32, height)
}

func TestImageDimensionsReturnsZeroForInvalidData(t *testing.T) {
	width, height := ImageDimensions(strings.NewReader("not an image"))

	assert.Equal(t, 0, width)
	assert.Equal(t, 0, height)
}
