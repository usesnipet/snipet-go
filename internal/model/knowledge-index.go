package model

import "github.com/usesnipet/snipet/pkg/jsonx"

type KnowledgeIndex struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string        `gorm:"type:varchar(255);not null" json:"name"`
	Driver        string        `gorm:"type:varchar(100);not null" json:"driver"`
	Configuration jsonx.JSONMap `gorm:"type:jsonb;not null" json:"configuration"`

	KnowledgeID string `gorm:"type:uuid;not null;index" json:"knowledge_id"`

	Knowledge *Knowledge             `gorm:"foreignKey:KnowledgeID;references:ID;constraint:OnDelete:CASCADE" json:"knowledge"`
	Items     []IndexedKnowledgeItem `gorm:"foreignKey:IndexID;constraint:OnDelete:CASCADE" json:"items"`
}
