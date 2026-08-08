package fileutil

import (
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

// MapKind classifies an already-normalized media type into a SourceItemKind.
func MapKind(mediaType string) knowledge.SourceItemKind {
	if mediaType == "" {
		return knowledge.SourceItemKindUnknown
	}

	switch mediaType {
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
