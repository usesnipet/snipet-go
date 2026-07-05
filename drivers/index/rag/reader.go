package rag

import (
	"context"

	"github.com/usesnipet/snipet/internal/runtime"
)

type Reader struct {
}

func NewReader(cfg Config) runtime.IIndexReader {
	return &Reader{}
}

func (r *Reader) Read(ctx context.Context, query string) ([]runtime.SourceItem, error) {
	return nil, nil
}
