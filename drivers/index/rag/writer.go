package rag

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/runtime"
)

type Writer struct {
	*components
}

func NewWriter(ctx context.Context, cfg Config) (runtime.IIndexWriter, error) {
	comps, err := newComponents(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Writer{components: comps}, nil
}

func (w *Writer) Index(ctx context.Context, itemID string, kind runtime.SourceItemKind, content any, attributes any) error {
	if itemID == "" {
		return fmt.Errorf("rag: indexed item id is required")
	}

	chunks, err := w.chunker.Chunk(kind, content, attributes, nil)
	if err != nil {
		return fmt.Errorf("rag: chunk: %w", err)
	}

	embeddings, err := w.embedder.Embed(ctx, chunks...)
	if err != nil {
		return fmt.Errorf("rag: embed: %w", err)
	}

	if err := w.store.Replace(ctx, itemID, embeddings); err != nil {
		return fmt.Errorf("rag: store: %w", err)
	}
	return nil
}

func (w *Writer) DeleteMany(ctx context.Context, ids []string) error {
	return w.store.DeleteMany(ctx, ids)
}

func (w *Writer) SupportedKinds() []runtime.SourceItemKind {
	return []runtime.SourceItemKind{
		runtime.SourceItemKindDocument,
		runtime.SourceItemKindText,
	}
}

func (w *Writer) Close() error {
	return w.components.Close()
}
