package model

import (
	"time"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

type KnowledgeItem struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ExternalID   string                   `gorm:"type:varchar(255);uniqueIndex:idx_knowledge_items_knowledge_external_id,priority:2" json:"external_id"`
	Name         string                   `gorm:"type:text" json:"name"`
	Hash         string                   `gorm:"type:varchar(128);index" json:"hash"`
	Metadata     util.JSONMap             `gorm:"type:jsonb" json:"metadata"`
	Attributes   util.JSONMap             `gorm:"type:jsonb" json:"attributes"`
	Kind         knowledge.SourceItemKind `gorm:"type:varchar(255)" json:"kind"`
	LastModified *time.Time            `json:"last_modified,omitempty"`

	KnowledgeID string `gorm:"type:uuid;not null;uniqueIndex:idx_knowledge_items_knowledge_external_id,priority:1" json:"knowledge_id"`

	Knowledge Knowledge              `gorm:"foreignKey:KnowledgeID" json:"-"`
	Indexes   []IndexedKnowledgeItem `gorm:"foreignKey:KnowledgeItemID" json:"-"`
}
