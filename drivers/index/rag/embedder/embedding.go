package embedder

import "github.com/usesnipet/snipet/drivers/index/rag/chunker"

type Embedding struct {
	Chunk     chunker.Chunk
	Embedding []float32
	Metadata  map[string]any
}

func NewEmbedding(chunk chunker.Chunk, embedding []float32, metadata map[string]any) *Embedding {
	return &Embedding{
		Chunk:     chunk,
		Embedding: embedding,
		Metadata:  metadata,
	}
}
