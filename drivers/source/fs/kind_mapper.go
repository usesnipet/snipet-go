package fs

import (
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/usesnipet/snipet/internal/runtime"
)

func detectByContent(item os.FileInfo) (*string, bool) {
	file, err := os.Open(item.Name())
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

func detectByExtension(item os.FileInfo) (*string, bool) {
	extension := path.Ext(item.Name())
	mimeType := mime.TypeByExtension(extension)
	return &mimeType, true
}

func normalizeMimeType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	return mediaType
}

func mapKind(item os.FileInfo) []runtime.SourceItemKind {
	mimeType, ok := detectByContent(item)
	if !ok {
		mimeType, ok = detectByExtension(item)
		if !ok {
			return []runtime.SourceItemKind{runtime.SourceItemKindUnknown}
		}
	}

	normalized := normalizeMimeType(*mimeType)

	switch normalized {
	case "text/plain":
		return []runtime.SourceItemKind{runtime.SourceItemKindText}
	case "application/pdf", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.ms-excel.sheet.binary.macroEnabled.12":
		return []runtime.SourceItemKind{runtime.SourceItemKindDocument}
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml":
		return []runtime.SourceItemKind{runtime.SourceItemKindImage}
	case "audio/mpeg", "audio/mp3", "audio/ogg", "audio/wav", "audio/x-wav":
		return []runtime.SourceItemKind{runtime.SourceItemKindAudio}
	case "video/mp4", "video/webm", "video/ogg", "video/quicktime":
		return []runtime.SourceItemKind{runtime.SourceItemKindVideo}
	case "application/json", "application/xml", "application/yaml", "application/toml", "application/hcl", "application/sql":
		return []runtime.SourceItemKind{runtime.SourceItemKindStructured}
	default:
		return []runtime.SourceItemKind{runtime.SourceItemKindUnknown}
	}
}
