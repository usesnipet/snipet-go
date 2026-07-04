package runtime

type SourceItemKind string

const (
	SourceItemKindText       SourceItemKind = "text"
	SourceItemKindDocument   SourceItemKind = "document"
	SourceItemKindImage      SourceItemKind = "image"
	SourceItemKindAudio      SourceItemKind = "audio"
	SourceItemKindVideo      SourceItemKind = "video"
	SourceItemKindStructured SourceItemKind = "structured"
	SourceItemKindUnknown    SourceItemKind = "unknown"
)
