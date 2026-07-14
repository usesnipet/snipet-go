package runtime

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

type IndexRecord struct {
	Content  any
	Metadata util.JSONMap
}

type IIndexReader interface {
	Read(ctx context.Context, query string, limit int) ([]IndexRecord, error)
	Close() error
}

type IIndexWriter interface {
	Index(ctx context.Context, itemID string, kind SourceItemKind, content any, attributes any) error
	DeleteMany(ctx context.Context, itemIDs []string) error
	SupportedKinds() []SourceItemKind
	Close() error
}

type IIndexDriver interface {
	Reader(config util.JSONMap) (IIndexReader, error)
	Writer(config util.JSONMap) (IIndexWriter, error)
	TestConnection(ctx context.Context, config util.JSONMap) error
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
}
