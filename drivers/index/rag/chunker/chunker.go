package chunker

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ledongthuc/pdf"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type Chunk struct {
	SeqID    int
	Content  string
	Start    int
	End      int
	Metadata map[string]any
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

func (c *Chunker) Chunk(kind knowledge.SourceItemKind, content, attributes any, metadata map[string]any) ([]Chunk, error) {
	text, err := extractText(kind, content, attributes)
	if err != nil {
		return nil, err
	}

	parts := splitTextParts(text, c.cfg)
	chunks := make([]Chunk, len(parts))
	for i, part := range parts {
		chunks[i] = Chunk{
			SeqID:    i,
			Content:  part.Content,
			Start:    part.Start,
			End:      part.End,
			Metadata: metadata,
		}
	}

	return chunks, nil
}

func extractText(kind knowledge.SourceItemKind, content, attributes any) (string, error) {
	switch kind {
	case knowledge.SourceItemKindText:
		text, ok := content.(string)
		if !ok {
			return "", fmt.Errorf("chunker: text content must be a string, got %T", content)
		}
		return text, nil
	case knowledge.SourceItemKindDocument:
		mediaType, err := documentMediaType(attributes)
		if err != nil {
			return "", err
		}
		if mediaType != "application/pdf" {
			return "", fmt.Errorf("chunker: unsupported document media type %q", mediaType)
		}
		return extractPDFText(content)
	default:
		return "", fmt.Errorf("chunker: unsupported content kind %q", kind)
	}
}

func documentMediaType(attributes any) (string, error) {
	switch v := attributes.(type) {
	case knowledge.DocumentAttributes:
		return v.MediaType, nil
	case *knowledge.DocumentAttributes:
		if v == nil {
			return "", fmt.Errorf("chunker: document attributes are nil")
		}
		return v.MediaType, nil
	case jsonx.JSONMap:
		attrs, err := jsonx.ParseJSONMap[knowledge.DocumentAttributes](v)
		if err != nil {
			return "", fmt.Errorf("chunker: parse document attributes: %w", err)
		}
		return attrs.MediaType, nil
	case map[string]any:
		attrs, err := jsonx.ParseJSONMap[knowledge.DocumentAttributes](jsonx.JSONMap(v))
		if err != nil {
			return "", fmt.Errorf("chunker: parse document attributes: %w", err)
		}
		return attrs.MediaType, nil
	default:
		return "", fmt.Errorf("chunker: unexpected document attributes type %T", attributes)
	}
}

func extractPDFText(content any) (string, error) {
	data, err := readContentBytes(content)
	if err != nil {
		return "", err
	}

	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("chunker: open pdf: %w", err)
	}

	plain, err := reader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("chunker: extract pdf text: %w", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(plain); err != nil {
		return "", fmt.Errorf("chunker: read pdf text: %w", err)
	}
	return buf.String(), nil
}

func readContentBytes(content any) ([]byte, error) {
	switch v := content.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case io.Reader:
		data, err := io.ReadAll(v)
		if err != nil {
			return nil, fmt.Errorf("chunker: read content: %w", err)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("chunker: pdf content must be []byte or io.Reader, got %T", content)
	}
}
