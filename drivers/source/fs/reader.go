package fs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type Reader struct {
	path string
	item *knowledge.SourceItem
}

func NewReader(config jsonx.JSONMap, itemID string) (knowledge.IKnowledgeReader, error) {
	cfg, err := jsonx.ParseJSONMap[Config](config)
	if err != nil {
		return nil, err
	}

	sourcePath := itemID
	if !filepath.IsAbs(itemID) {
		sourcePath = filepath.Join(cfg.BasePath, itemID)
	}

	sourceItem, _, err := sourceItemFromFile(sourcePath)
	if err != nil {
		return nil, err
	}
	return &Reader{
		path: sourcePath,
		item: sourceItem,
	}, nil
}

func (r *Reader) Kind() knowledge.SourceItemKind {
	return r.item.Kind
}

func (r *Reader) Open(ctx context.Context) (any, error) {
	switch r.item.Kind {
	case knowledge.SourceItemKindText:
		data, err := os.ReadFile(r.path)
		if err != nil {
			return nil, fmt.Errorf("fs: read text: %w", err)
		}
		return string(data), nil
	case knowledge.SourceItemKindDocument:
		data, err := os.ReadFile(r.path)
		if err != nil {
			return nil, fmt.Errorf("fs: read document: %w", err)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("fs: open not supported for kind %q", r.item.Kind)
	}
}

func (r *Reader) Attributes() any {
	return r.item.Attributes
}

func (r *Reader) Close() error {
	return nil
}
