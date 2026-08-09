package knowledge

// DocumentAttributes describes a SourceItemKindDocument item (e.g. PDF, DOCX):
// its media type, size in bytes, and best-effort bibliographic metadata.
type DocumentAttributes struct {
	MediaType string  `json:"media_type" validate:"required"`
	Size      int64   `json:"size" validate:"required"`
	Language  *string `json:"language" validate:"omitempty,len=2"`
	Title     *string `json:"title" validate:"omitempty,max=255"`
	Author    *string `json:"author" validate:"omitempty,max=255"`
}

// TextAttributes describes a SourceItemKindText item: plain text content
// with optional bibliographic metadata.
type TextAttributes struct {
	MediaType string  `json:"media_type" validate:"required"`
	Size      int64   `json:"size" validate:"required"`
	Language  *string `json:"language" validate:"omitempty,len=2"`
	Title     *string `json:"title" validate:"omitempty,max=255"`
	Author    *string `json:"author" validate:"omitempty,max=255"`
}

// ImageAttributes describes a SourceItemKindImage item: media type, size,
// and pixel dimensions.
type ImageAttributes struct {
	MediaType string `json:"media_type" validate:"required"`
	Size      int64  `json:"size" validate:"required"`
	Width     int    `json:"width" validate:"required"`
	Height    int    `json:"height" validate:"required"`
}

// AudioAttributes describes a SourceItemKindAudio item: media type, size,
// and playback characteristics that may not always be known.
type AudioAttributes struct {
	MediaType  string `json:"media_type" validate:"required"`
	Size       int64  `json:"size" validate:"required"`
	Duration   *int   `json:"duration" validate:"omitempty"`
	Bitrate    *int   `json:"bitrate" validate:"omitempty"`
	SampleRate *int   `json:"sample_rate" validate:"omitempty"`
	Channels   *int   `json:"channels" validate:"omitempty"`
}

// VideoAttributes describes a SourceItemKindVideo item: media type, size,
// and playback/frame characteristics that may not always be known.
type VideoAttributes struct {
	MediaType string `json:"media_type" validate:"required"`
	Size      int64  `json:"size" validate:"required"`
	Duration  *int   `json:"duration" validate:"omitempty"`
	Bitrate   *int   `json:"bitrate" validate:"omitempty"`
	FrameRate *int   `json:"frame_rate" validate:"omitempty"`
	Width     *int   `json:"width" validate:"omitempty"`
	Height    *int   `json:"height" validate:"omitempty"`
}

// StructuredAttributes describes a SourceItemKindStructured item (e.g. JSON,
// CSV): just its media type and size, since structure is implied by content.
type StructuredAttributes struct {
	MediaType string `json:"media_type" validate:"required"`
	Size      int64  `json:"size" validate:"required"`
}
