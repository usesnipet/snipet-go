package model

import (
	"time"

	"github.com/usesnipet/snipet/internal/util"
)

type IndexStatus string

const (
	IndexedStatusPending   IndexStatus = "pending"
	IndexedStatusSyncing   IndexStatus = "syncing"
	IndexedStatusIndexed   IndexStatus = "indexed"
	IndexedStatusProcessed IndexStatus = "skipped"
	IndexedStatusError     IndexStatus = "error"
)

type IndexedKnowledgeItem struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Hash      string       `gorm:"type:varchar(128);index" json:"hash"`
	IndexedAt *time.Time   `json:"indexed_at,omitempty"`
	Metadata  util.JSONMap `gorm:"type:jsonb" json:"metadata"`
	Status    IndexStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Reason    string       `gorm:"type:text" json:"reason"`
	LastError string       `gorm:"type:text" json:"last_error"`

	IndexID         string `gorm:"type:uuid;not null;index" json:"index_id"`
	KnowledgeItemID string `gorm:"type:uuid;not null;index" json:"knowledge_item_id"`

	Index         KnowledgeIndex `gorm:"foreignKey:IndexID" json:"-"`
	KnowledgeItem KnowledgeItem  `gorm:"foreignKey:KnowledgeItemID" json:"-"`
}
