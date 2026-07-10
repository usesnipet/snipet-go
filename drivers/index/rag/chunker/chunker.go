package chunker

import (
	"fmt"
	"io"

	"github.com/usesnipet/snipet/internal/runtime"
)

type Chunk struct {
	SeqID   int
	Content string
	Start   int
	End     int
}

type Chunker struct {
	cfg Config
}

func New(cfg Config) (*Chunker, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Chunker{cfg: cfg}, nil
}

func (c *Chunker) Chunk(content runtime.IContent) ([]Chunk, error) {
	text, err := extractText(content)
	if err != nil {
		return nil, err
	}

	parts := splitTextParts(text, c.cfg)
	chunks := make([]Chunk, len(parts))
	for i, part := range parts {
		chunks[i] = Chunk{
			SeqID:   i,
			Content: part.Content,
			Start:   part.Start,
			End:     part.End,
		}
	}

	return chunks, nil
}

func extractText(content runtime.IContent) (string, error) {
	switch value := content.(type) {
	case *runtime.TextContent:
		return value.Text, nil
	case *runtime.DocumentContent:
		if value.Doc == nil {
			return "", fmt.Errorf("chunker: document content has no reader")
		}
		data, err := io.ReadAll(value.Doc)
		if err != nil {
			return "", fmt.Errorf("chunker: read document: %w", err)
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("chunker: unsupported content kind %q", content.Kind())
	}
}
