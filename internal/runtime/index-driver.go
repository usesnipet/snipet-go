package runtime

import "context"

type IIndexDriver interface {
	Index(ctx context.Context, item SourceItem) error
	Search(ctx context.Context, query string) ([]SourceItem, error)
	Delete(ctx context.Context, item SourceItem) error
}
