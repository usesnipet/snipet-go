package rag

import (
	"context"

	"github.com/usesnipet/snipet/internal/runtime"
)

type Writer struct {
}

func NewWriter(cfg Config) (runtime.IIndexWriter, error) {
	return &Writer{}, nil
}

func (w *Writer) Index(ctx context.Context, item runtime.IContent) error {

	return nil
}

func (w *Writer) DeleteMany(ctx context.Context, ids []string) error {

	return nil
}

func (w *Writer) SupportedKinds() []runtime.SourceItemKind {
	return []runtime.SourceItemKind{
		runtime.SourceItemKindDocument,
		runtime.SourceItemKindText,
	}
}
