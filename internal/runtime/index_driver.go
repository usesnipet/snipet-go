package runtime

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

type IIndexReader interface {
	Read(ctx context.Context, query string) ([]SourceItem, error)
}

type IIndexWriter interface {
	Index(ctx context.Context, item SourceItem) error
	Delete(ctx context.Context, item SourceItem) error
	CanIndex(ctx context.Context, item SourceItem) (bool, error)
}

type IIndexDriver interface {
	Reader(config util.JSONMap) IIndexReader
	Writer(config util.JSONMap) IIndexWriter
	TestConnection(ctx context.Context, config util.JSONMap) error
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
}
