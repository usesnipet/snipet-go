package model

import (
	"time"

	"github.com/usesnipet/snipet/internal/util"
)

type KnowledgeType string

const (
	KnowledgeTypeDocuments KnowledgeType = "documents"
)

type SyncStatus string

const (
	SyncStatusInProgress SyncStatus = "in_progress"
	SyncStatusFailed     SyncStatus = "failed"
	SyncStatusSuccess    SyncStatus = "success"
)

type Knowledge struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string        `gorm:"type:varchar(255);not null" json:"name"`
	Description   string        `gorm:"type:text" json:"description"`
	Type          KnowledgeType `gorm:"type:varchar(100);not null;index" json:"type"`
	Driver        string        `gorm:"type:varchar(100);not null" json:"driver"`
	Configuration util.JSONMap  `gorm:"type:jsonb;not null" json:"configuration"`
	LastSyncedAt  *time.Time    `gorm:"type:timestamp;default:null" json:"last_synced_at"`
	SyncStatus    SyncStatus    `gorm:"type:varchar(20);default:null" json:"sync_status"`
	SyncError     *string       `gorm:"type:text;default:null" json:"sync_error"`

	Items   []KnowledgeItem  `gorm:"foreignKey:KnowledgeID;constraint:OnDelete:CASCADE" json:"-"`
	Indexes []KnowledgeIndex `gorm:"foreignKey:KnowledgeID;constraint:OnDelete:CASCADE" json:"-"`
}
