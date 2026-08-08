package fileutil

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// oneByteReader returns at most one byte per Read call, simulating a slow
// network stream where a single Read never fills the caller's buffer.
type oneByteReader struct {
	data []byte
	pos  int
}

func (r *oneByteReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	p[0] = r.data[r.pos]
	r.pos++
	return 1, nil
}

func TestDetectMediaTypeByContentDetectsPNG(t *testing.T) {
	pngHeader := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	body := append(pngHeader, bytes.Repeat([]byte{0}, 100)...)

	mimeType, ok := DetectMediaTypeByContent(bytes.NewReader(body))

	assert.True(t, ok)
	assert.Equal(t, "image/png", mimeType)
}

func TestDetectMediaTypeByContentHandlesSmallFile(t *testing.T) {
	mimeType, ok := DetectMediaTypeByContent(strings.NewReader("hi"))

	assert.True(t, ok)
	assert.Equal(t, "text/plain; charset=utf-8", mimeType)
}

func TestDetectMediaTypeByContentHandlesEmptyReader(t *testing.T) {
	mimeType, ok := DetectMediaTypeByContent(strings.NewReader(""))

	assert.False(t, ok)
	assert.Empty(t, mimeType)
}

func TestDetectMediaTypeByContentSurvivesPartialReads(t *testing.T) {
	pngHeader := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	body := append(pngHeader, bytes.Repeat([]byte{0}, 600)...)

	mimeType, ok := DetectMediaTypeByContent(&oneByteReader{data: body})

	assert.True(t, ok)
	assert.Equal(t, "image/png", mimeType)
}

func TestDetectMediaTypeByExtensionKnownExtension(t *testing.T) {
	mimeType, ok := DetectMediaTypeByExtension("report.json")

	assert.True(t, ok)
	assert.Contains(t, mimeType, "application/json")
}

func TestDetectMediaTypeByExtensionUnknownExtension(t *testing.T) {
	mimeType, ok := DetectMediaTypeByExtension("file.unknownext")

	assert.False(t, ok)
	assert.Empty(t, mimeType)
}

func TestNormalizeMediaTypeStripsParameters(t *testing.T) {
	normalized := NormalizeMediaType("text/plain; charset=utf-8")

	assert.Equal(t, "text/plain", normalized)
}

func TestNormalizeMediaTypeReturnsInputOnParseError(t *testing.T) {
	normalized := NormalizeMediaType("not a mime type;;;")

	assert.Equal(t, "not a mime type;;;", normalized)
}

func TestDetectMediaTypeReturnsNormalizedContentType(t *testing.T) {
	mediaType := DetectMediaType("file.txt", strings.NewReader("hello"))

	assert.Equal(t, "text/plain", mediaType)
}

func TestDetectMediaTypeFallsBackToExtensionOnEmptyReader(t *testing.T) {
	mediaType := DetectMediaType("data.json", strings.NewReader(""))

	assert.Equal(t, "application/json", mediaType)
}

func TestDetectMediaTypeReturnsEmptyWhenNothingMatches(t *testing.T) {
	mediaType := DetectMediaType("data.unknownext", strings.NewReader(""))

	assert.Empty(t, mediaType)
}
