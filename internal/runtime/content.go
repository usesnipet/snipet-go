package runtime

import "io"

type IContent interface {
	Kind() SourceItemKind
}

type DocumentContent struct {
	Doc io.ReadCloser
}

func (d *DocumentContent) Kind() SourceItemKind {
	return SourceItemKindDocument
}

type TextContent struct {
	Text string
}

func (t *TextContent) Kind() SourceItemKind {
	return SourceItemKindText
}

type ImageContent struct{}

func (i *ImageContent) Kind() SourceItemKind {
	return SourceItemKindImage
}

type AudioContent struct{}

func (a *AudioContent) Kind() SourceItemKind {
	return SourceItemKindAudio
}

type VideoContent struct{}

func (v *VideoContent) Kind() SourceItemKind {
	return SourceItemKindVideo
}

type StructuredContent struct{}

func (s *StructuredContent) Kind() SourceItemKind {
	return SourceItemKindStructured
}

type UnknownContent struct{}

func (u *UnknownContent) Kind() SourceItemKind {
	return SourceItemKindUnknown
}
