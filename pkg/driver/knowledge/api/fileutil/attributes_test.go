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
)

func TestNewDocumentAttributesFillsAllFields(t *testing.T) {
	attrs := NewDocumentAttributes("application/pdf", 1024, new("en"), new("Report"), new("Ada"))

	assert.Equal(t, "application/pdf", attrs.MediaType)
	assert.Equal(t, int64(1024), attrs.Size)
	assert.Equal(t, new("en"), attrs.Language)
	assert.Equal(t, new("Report"), attrs.Title)
	assert.Equal(t, new("Ada"), attrs.Author)
}

func TestNewDocumentAttributesAllowsNilMetadata(t *testing.T) {
	attrs := NewDocumentAttributes("application/pdf", 1024, nil, nil, nil)

	assert.Nil(t, attrs.Language)
	assert.Nil(t, attrs.Title)
	assert.Nil(t, attrs.Author)
}

func TestNewTextAttributesFillsAllFields(t *testing.T) {
	attrs := NewTextAttributes("text/plain", 1024, new("pt"), new("Notes"), new("Mayron"))

	assert.Equal(t, "text/plain", attrs.MediaType)
	assert.Equal(t, int64(1024), attrs.Size)
	assert.Equal(t, new("pt"), attrs.Language)
	assert.Equal(t, new("Notes"), attrs.Title)
	assert.Equal(t, new("Mayron"), attrs.Author)
}

func TestNewImageAttributesFillsAllFields(t *testing.T) {
	attrs := NewImageAttributes("image/png", 2048, 100, 200)

	assert.Equal(t, "image/png", attrs.MediaType)
	assert.Equal(t, int64(2048), attrs.Size)
	assert.Equal(t, 100, attrs.Width)
	assert.Equal(t, 200, attrs.Height)
}

func TestNewAudioAttributesFillsAllFields(t *testing.T) {
	attrs := NewAudioAttributes("audio/mpeg", 512, new(180), new(320), new(44100), new(2))

	assert.Equal(t, "audio/mpeg", attrs.MediaType)
	assert.Equal(t, int64(512), attrs.Size)
	assert.Equal(t, new(180), attrs.Duration)
	assert.Equal(t, new(320), attrs.Bitrate)
	assert.Equal(t, new(44100), attrs.SampleRate)
	assert.Equal(t, new(2), attrs.Channels)
}

func TestNewAudioAttributesAllowsNilMetadata(t *testing.T) {
	attrs := NewAudioAttributes("audio/mpeg", 512, nil, nil, nil, nil)

	assert.Nil(t, attrs.Duration)
	assert.Nil(t, attrs.Bitrate)
	assert.Nil(t, attrs.SampleRate)
	assert.Nil(t, attrs.Channels)
}

func TestNewVideoAttributesFillsAllFields(t *testing.T) {
	attrs := NewVideoAttributes("video/mp4", 4096, new(60), new(5000), new(30), new(1920), new(1080))

	assert.Equal(t, "video/mp4", attrs.MediaType)
	assert.Equal(t, int64(4096), attrs.Size)
	assert.Equal(t, new(60), attrs.Duration)
	assert.Equal(t, new(5000), attrs.Bitrate)
	assert.Equal(t, new(30), attrs.FrameRate)
	assert.Equal(t, new(1920), attrs.Width)
	assert.Equal(t, new(1080), attrs.Height)
}

func TestNewVideoAttributesAllowsNilMetadata(t *testing.T) {
	attrs := NewVideoAttributes("video/mp4", 4096, nil, nil, nil, nil, nil)

	assert.Nil(t, attrs.Duration)
	assert.Nil(t, attrs.Bitrate)
	assert.Nil(t, attrs.FrameRate)
	assert.Nil(t, attrs.Width)
	assert.Nil(t, attrs.Height)
}

func TestNewStructuredAttributesFillsAllFields(t *testing.T) {
	attrs := NewStructuredAttributes("application/json", 64)

	assert.Equal(t, "application/json", attrs.MediaType)
	assert.Equal(t, int64(64), attrs.Size)
}

func TestToJSONMapSerializesTypedAttributes(t *testing.T) {
	m := ToJSONMap(NewAudioAttributes("audio/mpeg", 512, new(180), nil, nil, nil))

	assert.Equal(t, "audio/mpeg", m["media_type"])
	assert.Equal(t, float64(512), m["size"])
	assert.Equal(t, float64(180), m["duration"])
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
