package model

import "github.com/usesnipet/snipet/internal/util"

type KnowledgeIndex struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string       `gorm:"type:varchar(255);not null" json:"name"`
	Driver        string       `gorm:"type:varchar(100);not null" json:"driver"`
	Configuration util.JSONMap `gorm:"type:jsonb;not null" json:"configuration"`

	KnowledgeID string `gorm:"type:uuid;not null;index" json:"knowledge_id"`

	Knowledge Knowledge              `gorm:"foreignKey:KnowledgeID" json:"-"`
	Items     []IndexedKnowledgeItem `gorm:"foreignKey:IndexID" json:"-"`
}
