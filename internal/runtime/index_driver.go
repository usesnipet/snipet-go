package runtime

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

type IIndexReader interface {
	Read(ctx context.Context, query string) ([]SourceItem, error)
}

type IIndexWriter interface {
	Index(ctx context.Context, content IContent) error
	DeleteMany(ctx context.Context, itemIds []string) error
	SupportedKinds() []SourceItemKind
}

type IIndexDriver interface {
	Reader(config util.JSONMap) (IIndexReader, error)
	Writer(config util.JSONMap) (IIndexWriter, error)
	TestConnection(ctx context.Context, config util.JSONMap) error
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
}
