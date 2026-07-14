package runtime

import (
	"context"
	"time"

	"github.com/usesnipet/snipet/internal/util"
)

type SourceItem struct {
	ID           string
	Name         string
	Kind         SourceItemKind
	Metadata     util.JSONMap
	Attributes   util.JSONMap
	LastModified *time.Time
}

type ISourceIterator interface {
	Next(ctx context.Context) bool
	Item() *SourceItem
	GetHash() string
	Err() error
	Close() error
}

type ISourceReader interface {
	Kind() SourceItemKind
	Attributes() any
	Open(ctx context.Context) (any, error)
	Close() error
}

type ISourceDriver interface {
	Iterator(ctx context.Context, config util.JSONMap) (ISourceIterator, error)
	Reader(ctx context.Context, config util.JSONMap, itemID string) (ISourceReader, error)
	TestConnection(ctx context.Context, config util.JSONMap) error
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
}
