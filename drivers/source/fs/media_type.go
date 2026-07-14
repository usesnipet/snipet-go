package fs

import (
	"mime"
	"net/http"
	"os"
	"path"
)

func detectByContent(filePath string) (*string, bool) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return nil, false
	}

	mimeType := http.DetectContentType(buffer[:n])
	return &mimeType, true
}

func detectByExtension(filePath string) (*string, bool) {
	extension := path.Ext(filePath)
	mimeType := mime.TypeByExtension(extension)
	if mimeType == "" {
		return nil, false
	}
	return &mimeType, true
}

func normalizeMimeType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	return mediaType
}

func detectMediaType(path string) string {
	mimeType, ok := detectByContent(path)
	if !ok {
		mimeType, ok = detectByExtension(path)
		if !ok {
			return ""
		}
	}
	return normalizeMimeType(*mimeType)
}
