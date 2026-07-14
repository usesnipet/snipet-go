package runtime

type DocumentAttributes struct {
	MediaType string  `json:"media_type" validate:"required"`
	Size      int64   `json:"size" validate:"required"`
	Language  *string `json:"language" validate:"omitempty,len=2"`
	Title     *string `json:"title" validate:"omitempty,max=255"`
	Author    *string `json:"author" validate:"omitempty,max=255"`
}

type TextAttributes struct {
	Language *string `json:"language" validate:"omitempty,len=2"`
	Title    *string `json:"title" validate:"omitempty,max=255"`
	Author   *string `json:"author" validate:"omitempty,max=255"`
}

type ImageAttributes struct {
	MediaType string `json:"media_type" validate:"required"`
	Size      int64  `json:"size" validate:"required"`
	Width     int    `json:"width" validate:"required"`
	Height    int    `json:"height" validate:"required"`
}

type AudioAttributes struct {
	MediaType  string `json:"media_type" validate:"required"`
	Size       int64  `json:"size" validate:"required"`
	Duration   *int   `json:"duration" validate:"omitempty"`
	Bitrate    *int   `json:"bitrate" validate:"omitempty"`
	SampleRate *int   `json:"sample_rate" validate:"omitempty"`
	Channels   *int   `json:"channels" validate:"omitempty"`
}

type VideoAttributes struct {
	MediaType string `json:"media_type" validate:"required"`
	Size      int64  `json:"size" validate:"required"`
	Duration  *int   `json:"duration" validate:"omitempty"`
	Bitrate   *int   `json:"bitrate" validate:"omitempty"`
	FrameRate *int   `json:"frame_rate" validate:"omitempty"`
	Width     *int   `json:"width" validate:"omitempty"`
	Height    *int   `json:"height" validate:"omitempty"`
}

type StructuredAttributes struct {
	MediaType string `json:"media_type" validate:"required"`
	Size      int64  `json:"size" validate:"required"`
}
