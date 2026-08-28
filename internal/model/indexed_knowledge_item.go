package model

import (
	"time"

	"github.com/usesnipet/snipet/pkg/jsonx"
)

type IndexStatus string

const (
	IndexedStatusPending IndexStatus = "pending"
	IndexedStatusSyncing IndexStatus = "syncing"
	IndexedStatusIndexed IndexStatus = "indexed"
	IndexedStatusSkipped IndexStatus = "skipped"
	IndexedStatusError   IndexStatus = "error"
)

type IndexedKnowledgeItem struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Hash      string        `gorm:"type:varchar(128);index" json:"hash"`
	IndexedAt *time.Time    `json:"indexed_at,omitempty"`
	Metadata  jsonx.JSONMap `gorm:"type:jsonb" json:"metadata"`
	Status    IndexStatus   `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Reason    *string       `gorm:"type:text" json:"reason,omitempty"`
	LastError *string       `gorm:"type:text" json:"last_error,omitempty"`

	IndexID         string  `gorm:"type:uuid;not null;index" json:"index_id"`
	KnowledgeItemID *string `gorm:"type:uuid;index" json:"knowledge_item_id,omitempty"`

	Index         *KnowledgeIndex `gorm:"foreignKey:IndexID;references:ID;constraint:OnDelete:CASCADE" json:"index"`
	KnowledgeItem *KnowledgeItem  `gorm:"foreignKey:KnowledgeItemID;references:ID;constraint:OnDelete:SET NULL" json:"knowledge_item"`
}
