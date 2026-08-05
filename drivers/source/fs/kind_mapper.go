package fs

import (
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func mapKind(filePath string) knowledge.SourceItemKind {
	normalized := detectMediaType(filePath)
	if normalized == "" {
		return knowledge.SourceItemKindUnknown
	}

	switch normalized {
	case "text/plain":
		return knowledge.SourceItemKindText
	case "application/pdf", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.ms-excel.sheet.binary.macroEnabled.12":
		return knowledge.SourceItemKindDocument
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml":
		return knowledge.SourceItemKindImage
	case "audio/mpeg", "audio/mp3", "audio/ogg", "audio/wav", "audio/x-wav":
		return knowledge.SourceItemKindAudio
	case "video/mp4", "video/webm", "video/ogg", "video/quicktime":
		return knowledge.SourceItemKindVideo
	case "application/json", "application/xml", "application/yaml", "application/toml", "application/hcl", "application/sql":
		return knowledge.SourceItemKindStructured
	default:
		return knowledge.SourceItemKindUnknown
	}
}
