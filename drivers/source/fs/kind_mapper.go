package fs

import (
	"github.com/usesnipet/snipet/internal/runtime"
)

func mapKind(filePath string) runtime.SourceItemKind {
	normalized := detectMediaType(filePath)
	if normalized == "" {
		return runtime.SourceItemKindUnknown
	}

	switch normalized {
	case "text/plain":
		return runtime.SourceItemKindText
	case "application/pdf", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.ms-excel.sheet.binary.macroEnabled.12":
		return runtime.SourceItemKindDocument
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml":
		return runtime.SourceItemKindImage
	case "audio/mpeg", "audio/mp3", "audio/ogg", "audio/wav", "audio/x-wav":
		return runtime.SourceItemKindAudio
	case "video/mp4", "video/webm", "video/ogg", "video/quicktime":
		return runtime.SourceItemKindVideo
	case "application/json", "application/xml", "application/yaml", "application/toml", "application/hcl", "application/sql":
		return runtime.SourceItemKindStructured
	default:
		return runtime.SourceItemKindUnknown
	}
}
