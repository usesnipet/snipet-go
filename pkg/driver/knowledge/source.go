package knowledge

import (
	"context"
	"time"

	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// SourceItemKind classifies the content of a SourceItem, determining which
// *Attributes type (see attributes.go) describes it and which index writers
// can ingest it (see IKnowledgeIndexWriter.SupportedKinds).
type SourceItemKind string

const (
	SourceItemKindText       SourceItemKind = "text"
	SourceItemKindDocument   SourceItemKind = "document"
	SourceItemKindImage      SourceItemKind = "image"
	SourceItemKindAudio      SourceItemKind = "audio"
	SourceItemKindVideo      SourceItemKind = "video"
	SourceItemKindStructured SourceItemKind = "structured"
	SourceItemKindUnknown    SourceItemKind = "unknown"
)

// SourceItem is a single unit of content discovered in a knowledge source
// (e.g. a file, page, or record), identified by ID and enriched with
// metadata/attributes for downstream indexing.
type SourceItem struct {
	ID           string
	Name         string
	Kind         SourceItemKind
	Metadata     jsonx.JSONMap
	Attributes   jsonx.JSONMap
	LastModified *time.Time
}

// IKnowledgeIterator walks the items exposed by a source, cursor-style: call
// Next until it returns false, reading Item after each successful Next. Err
// reports any error that stopped iteration; Close releases underlying
// resources regardless of how iteration ended.
type IKnowledgeIterator interface {
	Next(ctx context.Context) bool
	Item() *SourceItem
	GetHash() string
	Err() error
	Close() error
}

// IKnowledgeReader opens the content of a single source item for reading.
// Kind and Attributes describe the item before Open is called; the concrete
// type returned by Open depends on Kind (e.g. an io.Reader for binary kinds).
type IKnowledgeReader interface {
	Kind() SourceItemKind
	Attributes() any
	Open(ctx context.Context) (any, error)
	Close() error
}

// ISourceDriver is the driver contract for a knowledge source backend (e.g.
// a filesystem, bucket, or third-party integration). Iterator enumerates the
// items available under a given config; Reader opens one item by ID.
type ISourceDriver interface {
	driver.IDriver

	Iterator(ctx context.Context, config jsonx.JSONMap) (IKnowledgeIterator, error)
	Reader(ctx context.Context, config jsonx.JSONMap, itemID string) (IKnowledgeReader, error)
}
