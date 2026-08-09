package fileutil

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func encodePNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := range width {
		for y := range height {
			img.Set(x, y, color.White)
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func TestInspectClassifiesImageAndBuildsAttributes(t *testing.T) {
	content := encodePNG(t, 64, 32)

	kind, attrs, _ := Inspect("photo.png", int64(len(content)), bytes.NewReader(content))

	assert.Equal(t, knowledge.SourceItemKindImage, kind)
	assert.Equal(t, "image/png", attrs["media_type"])
	assert.Equal(t, float64(64), attrs["width"])
	assert.Equal(t, float64(32), attrs["height"])
}

func TestInspectClassifiesStructuredContent(t *testing.T) {
	content := []byte(`{"key": "value"}`)

	kind, attrs, _ := Inspect("data", int64(len(content)), bytes.NewReader(content))

	assert.Equal(t, knowledge.SourceItemKindStructured, kind)
	assert.Equal(t, "application/json", attrs["media_type"])
}

func TestInspectFallsBackToUnknownForUnrecognizedContent(t *testing.T) {
	// Control bytes with no valid text shape and no known magic signature.
	content := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}

	kind, attrs, _ := Inspect("data.unknownext", int64(len(content)), bytes.NewReader(content))

	assert.Equal(t, knowledge.SourceItemKindUnknown, kind)
	assert.Empty(t, attrs)
}

func TestInspectReaderReplaysFullContentForSmallInput(t *testing.T) {
	content := []byte("hello world")

	_, _, rest := Inspect("file.txt", int64(len(content)), bytes.NewReader(content))

	replayed, err := io.ReadAll(rest)
	require.NoError(t, err)
	assert.Equal(t, content, replayed)
}

func TestInspectReaderReplaysFullContentLargerThanSniffBuffer(t *testing.T) {
	content := make([]byte, contentSniffLength*3)
	for i := range content {
		content[i] = byte(i % 251)
	}

	_, _, rest := Inspect("data.bin", int64(len(content)), bytes.NewReader(content))

	replayed, err := io.ReadAll(rest)
	require.NoError(t, err)
	assert.Equal(t, content, replayed)
}

func TestInspectReaderHashMatchesOriginalContent(t *testing.T) {
	content := []byte(strings.Repeat("snipet fileutil inspect ", 500))

	_, _, rest := Inspect("notes.txt", int64(len(content)), bytes.NewReader(content))

	gotHash, err := Hash(rest)
	require.NoError(t, err)

	wantHash, err := Hash(bytes.NewReader(content))
	require.NoError(t, err)

	assert.Equal(t, wantHash, gotHash)
}

func TestInspectReaderReplaysContentEvenWhenKindUnknown(t *testing.T) {
	content := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	_, _, rest := Inspect("data.unknownext", int64(len(content)), bytes.NewReader(content))

	replayed, err := io.ReadAll(rest)
	require.NoError(t, err)
	assert.Equal(t, content, replayed)
}
