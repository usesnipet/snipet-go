package knowledge

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
)

type KnowledgeIndexRecord struct {
	Content  any
	Metadata util.JSONMap
}

type IKnowledgeIndexReader interface {
	Read(ctx context.Context, query string, limit int) ([]KnowledgeIndexRecord, error)
	Close() error
}

type IKnowledgeIndexWriter interface {
	Index(ctx context.Context, itemID string, kind SourceItemKind, content any, attributes any) error
	DeleteMany(ctx context.Context, itemIDs []string) error
	SupportedKinds() []SourceItemKind
	Close() error
}

type IIndexDriver interface {
	driver.IDriver

	Reader(config util.JSONMap) (IKnowledgeIndexReader, error)
	Writer(config util.JSONMap) (IKnowledgeIndexWriter, error)
}
