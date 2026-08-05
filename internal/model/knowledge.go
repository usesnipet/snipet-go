package model

import (
	"time"

	"github.com/usesnipet/snipet/pkg/jsonx"
)

type SyncStatus string

const (
	SyncStatusPending    SyncStatus = "pending"
	SyncStatusInProgress SyncStatus = "in_progress"
	SyncStatusFailed     SyncStatus = "failed"
	SyncStatusSuccess    SyncStatus = "success"
)

type Knowledge struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string        `gorm:"type:varchar(255);not null" json:"name"`
	Description   string        `gorm:"type:text" json:"description"`
	Driver        string        `gorm:"type:varchar(100);not null" json:"driver"`
	Configuration jsonx.JSONMap `gorm:"type:jsonb;not null" json:"configuration"`
	LastSyncedAt  *time.Time    `gorm:"type:timestamp;default:null" json:"last_synced_at"`
	SyncStatus    SyncStatus    `gorm:"type:varchar(20);default:null" json:"sync_status"`
	SyncError     *string       `gorm:"type:text;default:null" json:"sync_error"`

	Items   []KnowledgeItem  `gorm:"foreignKey:KnowledgeID;constraint:OnDelete:CASCADE" json:"-"`
	Indexes []KnowledgeIndex `gorm:"foreignKey:KnowledgeID;constraint:OnDelete:CASCADE" json:"-"`
}
