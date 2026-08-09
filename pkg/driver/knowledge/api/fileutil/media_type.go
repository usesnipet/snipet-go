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

// readHeader reads up to contentSniffLength bytes from r. A reader shorter
// than that is not an error — only a genuine I/O failure is.
func readHeader(r io.Reader) ([]byte, error) {
	buffer := make([]byte, contentSniffLength)
	n, err := io.ReadFull(r, buffer)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}
	return buffer[:n], nil
}

func mimeFromHeader(header []byte) (string, bool) {
	if len(header) == 0 {
		return "", false
	}
	return mimetype.Detect(header).String(), true
}

// DetectMediaTypeByContent sniffs r's leading bytes for a magic-number or
// text-shape match (images, audio, video, archives, JSON/XML/HTML/plain
// text, etc). Unlike a plain single Read, it uses io.ReadFull so partial
// reads from non-local readers (e.g. network streams) don't cause a false
// "undetected" result.
func DetectMediaTypeByContent(r io.Reader) (string, bool) {
	header, err := readHeader(r)
	if err != nil {
		return "", false
	}
	return mimeFromHeader(header)
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

// mediaTypeFromHeader resolves the normalized media type for name from an
// already-read header, preferring content sniffing and falling back to the
// file extension.
func mediaTypeFromHeader(name string, header []byte) string {
	mimeType, ok := mimeFromHeader(header)
	if !ok {
		mimeType, ok = DetectMediaTypeByExtension(name)
		if !ok {
			return ""
		}
	}
	return NormalizeMediaType(mimeType)
}

// DetectMediaType resolves the normalized media type for name, preferring
// content sniffing from r and falling back to the file extension.
func DetectMediaType(name string, r io.Reader) string {
	header, err := readHeader(r)
	if err != nil {
		return ""
	}
	return mediaTypeFromHeader(name, header)
}
