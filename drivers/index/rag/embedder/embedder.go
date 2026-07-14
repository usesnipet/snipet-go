package embedder

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/drivers/index/rag/chunker"
	"google.golang.org/genai"
)

type Embedder struct {
	cfg     Config
	client  *genai.Client
	started bool
}

func New(cfg Config) (*Embedder, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Embedder{cfg: cfg}, nil
}

func (e *Embedder) Start(ctx context.Context) error {
	if e.started {
		return nil
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  e.cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}
	e.client = client
	e.started = true
	return nil
}

func (e *Embedder) Embed(ctx context.Context, chunks ...chunker.Chunk) ([]*Embedding, error) {
	if !e.started {
		return nil, fmt.Errorf("embedder: not started")
	}
	if len(chunks) == 0 {
		return nil, nil
	}

	contents := make([]*genai.Content, len(chunks))
	for i, chunk := range chunks {
		contents[i] = genai.NewContentFromText(chunk.Content, genai.RoleUser)
	}

	result, err := e.client.Models.EmbedContent(ctx, e.cfg.Model, contents, nil)
	if err != nil {
		return nil, err
	}
	if len(result.Embeddings) != len(chunks) {
		return nil, fmt.Errorf("embedder: expected %d embeddings, got %d", len(chunks), len(result.Embeddings))
	}

	embeddings := make([]*Embedding, len(chunks))
	for i, embedding := range result.Embeddings {
		metadata := chunks[i].Metadata
		if metadata == nil {
			metadata = map[string]any{}
		}
		embeddings[i] = NewEmbedding(chunks[i], embedding.Values, metadata)
	}
	return embeddings, nil
}

func (e *Embedder) EmbedQuery(ctx context.Context, query string) ([]float32, error) {
	embeddings, err := e.Embed(ctx, chunker.Chunk{Content: query})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 || embeddings[0] == nil {
		return nil, fmt.Errorf("embedder: empty query embedding")
	}
	return embeddings[0].Embedding, nil
}
