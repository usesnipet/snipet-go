package model

import "github.com/usesnipet/snipet/internal/util"

type KnowledgeType string

const (
	KnowledgeTypeDocuments KnowledgeType = "documents"
)

type Knowledge struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string        `gorm:"type:varchar(255);not null" json:"name"`
	Description   string        `gorm:"type:text" json:"description"`
	Type          KnowledgeType `gorm:"type:varchar(100);not null;index" json:"type"`
	Provider      string        `gorm:"type:varchar(100);not null" json:"provider"`
	Configuration util.JSONMap  `gorm:"type:jsonb;not null" json:"configuration"`

	Items   []KnowledgeItem  `gorm:"foreignKey:KnowledgeID;constraint:OnDelete:CASCADE" json:"-"`
	Indexes []KnowledgeIndex `gorm:"foreignKey:KnowledgeID;constraint:OnDelete:CASCADE" json:"-"`
}
