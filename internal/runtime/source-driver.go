package runtime

import (
	"context"
	"time"

	"github.com/usesnipet/snipet/internal/util"
)

type SourceItem struct {
	ID           string
	Name         string
	Metadata     util.JSONMap
	LastModified *time.Time
}

type SourceIterator interface {
	Next(ctx context.Context) bool
	Item() *SourceItem
	GetHash() string
	Err() error
	Close() error
}

type SourceDriver interface {
	Iterator(ctx context.Context, config util.JSONMap) (SourceIterator, error)
	TestConnection(ctx context.Context, config util.JSONMap) error
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
}
