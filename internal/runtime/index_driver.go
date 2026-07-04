package runtime

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

type IIndexDriver interface {
	ProcessableSourceTypes() []SourceType
	Index(ctx context.Context, item SourceItem) error
	Search(ctx context.Context, query string) ([]SourceItem, error)
	Delete(ctx context.Context, item SourceItem) error
	TestConnection(ctx context.Context, config util.JSONMap) error
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
}
