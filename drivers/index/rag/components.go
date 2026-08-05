package rag

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/drivers/index/rag/chunker"
	"github.com/usesnipet/snipet/drivers/index/rag/embedder"
	"github.com/usesnipet/snipet/drivers/index/rag/store"
)

type components struct {
	chunker  *chunker.Chunker
	embedder *embedder.Embedder
	store    *store.Store
}

func newComponents(ctx context.Context, cfg Config) (*components, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c, err := chunker.New(cfg.ChunkerConfig())
	if err != nil {
		return nil, err
	}

	e, err := embedder.New(cfg.EmbedderConfig())
	if err != nil {
		return nil, err
	}
	if err := e.Start(ctx); err != nil {
		return nil, fmt.Errorf("rag: start embedder: %w", err)
	}

	s, err := store.NewStore(cfg.StoreConfig())
	if err != nil {
		return nil, err
	}
	if err := s.Start(ctx); err != nil {
		return nil, fmt.Errorf("rag: start store: %w", err)
	}

	return &components{
		chunker:  c,
		embedder: e,
		store:    s,
	}, nil
}

func (c *components) Close() error {
	if c == nil {
		return nil
	}
	if c.store != nil {
		c.store.Close()
	}
	return nil
}
