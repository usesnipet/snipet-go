package rag

import (
	"context"

	"github.com/teilomillet/raggo"
	"github.com/usesnipet/snipet/internal/runtime"
)

type Writer struct {
	parser   *raggo.Parser
	chunker  *raggo.Chunker
	embedder *raggo.Embedder
	vectorDB *raggo.VectorDB
}

func NewWriter(cfg Config) (runtime.IIndexWriter, error) {
	parser := raggo.NewParser()
	chunker, err := raggo.NewChunker(raggo.ChunkSize(cfg.ChunkSize))
	if err != nil {
		return nil, err
	}

	embedder, err := raggo.NewEmbedder(
		raggo.SetEmbedderProvider(cfg.Embedder.Provider),
		raggo.SetEmbedderModel(cfg.Embedder.Model),
		raggo.SetEmbedderAPIKey(cfg.Embedder.APIKey),
	)
	if err != nil {
		return nil, err
	}

	vectorDB, err := raggo.NewVectorDB(
		raggo.WithType("milvus"),
		raggo.WithAddress(cfg.Milvus.Address),
		raggo.WithDimension(cfg.Milvus.Dimension),
		raggo.WithMaxPoolSize(cfg.Milvus.MaxPoolSize),
	)
	if err != nil {
		return nil, err
	}

	return &Writer{
		parser:   &parser,
		chunker:  &chunker,
		embedder: &embedder,
		vectorDB: vectorDB,
	}, nil
}

func (w *Writer) Index(ctx context.Context, item runtime.SourceItem) error {

	return nil
}

func (w *Writer) Delete(ctx context.Context, item runtime.SourceItem) error {
	return nil
}

func (w *Writer) CanIndex(ctx context.Context, item runtime.SourceItem) (bool, error) {
	return true, nil
}
