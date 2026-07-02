package model

import "github.com/usesnipet/snipet/internal/util"

type Knowledge struct {
	ID string `json:"id"`

	Name          string       `gorm:"type:varchar(255);not null" json:"title"`
	Description   string       `gorm:"type:text;not null" json:"description"`
	Provider      string       `gorm:"type:varchar(255);not null" json:"provider"`
	Configuration util.JSONMap `gorm:"type:jsonb;not null" json:"configuration"`

	KnowledgeToBots []BotToKnowledge `gorm:"foreignKey:KnowledgeID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
