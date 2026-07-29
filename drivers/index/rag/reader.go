package rag

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/internal/util"
)

type Reader struct {
	*components
}

func NewReader(ctx context.Context, cfg Config) (knowledge.IKnowledgeIndexReader, error) {
	comps, err := newComponents(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Reader{components: comps}, nil
}

func (r *Reader) Read(ctx context.Context, query string, limit int) ([]knowledge.KnowledgeIndexRecord, error) {
	if query == "" {
		return nil, fmt.Errorf("rag: query is required")
	}

	vector, err := r.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("rag: embed query: %w", err)
	}

	results, err := r.store.Search(ctx, vector, limit)
	if err != nil {
		return nil, fmt.Errorf("rag: search: %w", err)
	}

	records := make([]knowledge.KnowledgeIndexRecord, 0, len(results))
	for _, result := range results {
		metadata := util.JSONMap{
			"indexed_item_id": result.IndexedItemID,
			"seq_id":          result.Chunk.SeqID,
			"start":           result.Chunk.Start,
			"end":             result.Chunk.End,
			"score":           1 - result.Distance,
		}
		for k, v := range result.Metadata {
			metadata[k] = v
		}
		records = append(records, knowledge.KnowledgeIndexRecord{
			Content:  result.Chunk.Content,
			Metadata: metadata,
		})
	}
	return records, nil
}

func (r *Reader) Close() error {
	return r.components.Close()
}
