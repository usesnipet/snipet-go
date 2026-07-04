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
	_, err = file.Read(buffer)
	if err != nil {
		return nil, false
	}

	mimeType := http.DetectContentType(buffer)
	return &mimeType, true
}

func detectByExtension(item os.FileInfo) (*string, bool) {
	extension := path.Ext(item.Name())
	mimeType := mime.TypeByExtension(extension)
	return &mimeType, true
}

func mapKind(item os.FileInfo) []runtime.SourceItemKind {
	mimeType, ok := detectByContent(item)
	if !ok {
		mimeType, ok = detectByExtension(item)
		if !ok {
			return []runtime.SourceItemKind{runtime.SourceItemKindUnknown}
		}
	}

	switch *mimeType {
	case "text/plain":
		return []runtime.SourceItemKind{runtime.SourceItemKindText}
	case "application/pdf":
		return []runtime.SourceItemKind{runtime.SourceItemKindDocument}
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return []runtime.SourceItemKind{runtime.SourceItemKindImage}
	case "audio/mpeg", "audio/mp3", "audio/ogg", "audio/wav":
		return []runtime.SourceItemKind{runtime.SourceItemKindAudio}
	case "video/mp4", "video/webm", "video/ogg":
		return []runtime.SourceItemKind{runtime.SourceItemKindVideo}
	case "application/json", "application/xml", "application/yaml", "application/toml", "application/hcl":
		return []runtime.SourceItemKind{runtime.SourceItemKindStructured}
	default:
		return []runtime.SourceItemKind{runtime.SourceItemKindUnknown}
	}
}
