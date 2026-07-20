package driver

import (
	"context"
	"time"

	"github.com/usesnipet/snipet/internal/util"
)

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

type SourceItem struct {
	ID           string
	Name         string
	Kind         SourceItemKind
	Metadata     util.JSONMap
	Attributes   util.JSONMap
	LastModified *time.Time
}

type IKnowledgeIterator interface {
	Next(ctx context.Context) bool
	Item() *SourceItem
	GetHash() string
	Err() error
	Close() error
}

type IKnowledgeReader interface {
	Kind() SourceItemKind
	Attributes() any
	Open(ctx context.Context) (any, error)
	Close() error
}

type IKnowledgeSource interface {
	IDriver
	Iterator(ctx context.Context, config util.JSONMap) (IKnowledgeIterator, error)
	Reader(ctx context.Context, config util.JSONMap, itemID string) (IKnowledgeReader, error)
}
