package fileutil

import (
	"io"
	"mime"
	"path"

	"github.com/gabriel-vasile/mimetype"
)

// contentSniffLength matches mimetype's own default read limit, balancing
// detection accuracy (some formats, e.g. docx, need more than a few bytes)
// against how much of a remote-backed reader (s3, gdrive) gets pulled just
// to classify a file.
const contentSniffLength = 4096

// DetectMediaTypeByContent sniffs r's leading bytes for a magic-number or
// text-shape match (images, audio, video, archives, JSON/XML/HTML/plain
// text, etc). Unlike a plain single Read, it uses io.ReadFull so partial
// reads from non-local readers (e.g. network streams) don't cause a false
// "undetected" result.
func DetectMediaTypeByContent(r io.Reader) (string, bool) {
	buffer := make([]byte, contentSniffLength)
	n, err := io.ReadFull(r, buffer)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", false
	}
	if n == 0 {
		return "", false
	}
	return mimetype.Detect(buffer[:n]).String(), true
}

// DetectMediaTypeByExtension resolves a MIME type from name's file extension.
func DetectMediaTypeByExtension(name string) (string, bool) {
	extension := path.Ext(name)
	mimeType := mime.TypeByExtension(extension)
	if mimeType == "" {
		return "", false
	}
	return mimeType, true
}

// NormalizeMediaType strips parameters (e.g. "; charset=utf-8") from a
// Content-Type-style string, returning it unchanged if it fails to parse.
func NormalizeMediaType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	return mediaType
}

// DetectMediaType resolves the normalized media type for name, preferring
// content sniffing from r and falling back to the file extension.
func DetectMediaType(name string, r io.Reader) string {
	mimeType, ok := DetectMediaTypeByContent(r)
	if !ok {
		mimeType, ok = DetectMediaTypeByExtension(name)
		if !ok {
			return ""
		}
	}
	return NormalizeMediaType(mimeType)
}
