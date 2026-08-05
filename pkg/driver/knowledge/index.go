package knowledge

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// KnowledgeIndexRecord is a single match returned from an index query: the
// indexed content plus whatever metadata the index stored alongside it.
type KnowledgeIndexRecord struct {
	Content  any
	Metadata jsonx.JSONMap
}

// IKnowledgeIndexReader queries an already-populated index (e.g. a vector or
// full-text search backend). Close releases any underlying connection.
type IKnowledgeIndexReader interface {
	Read(ctx context.Context, query string, limit int) ([]KnowledgeIndexRecord, error)
	Close() error
}

// IKnowledgeIndexWriter ingests source items into an index. SupportedKinds
// reports which SourceItemKind values Index can accept; callers should skip
// items whose kind is not supported rather than calling Index with them.
type IKnowledgeIndexWriter interface {
	Index(ctx context.Context, itemID string, kind SourceItemKind, content any, attributes any) error
	DeleteMany(ctx context.Context, itemIDs []string) error
	SupportedKinds() []SourceItemKind
	Close() error
}

// IIndexDriver is the driver contract for a knowledge index backend. Reader
// and Writer are constructed lazily from a config map so a single driver
// instance can serve multiple index configurations.
type IIndexDriver interface {
	driver.IDriver

	Reader(config jsonx.JSONMap) (IKnowledgeIndexReader, error)
	Writer(config jsonx.JSONMap) (IKnowledgeIndexWriter, error)
}
