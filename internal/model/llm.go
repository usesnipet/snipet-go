package model

import (
	"github.com/usesnipet/snipet/internal/util"
)

type LLM struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string       `gorm:"type:varchar(255);not null" json:"name"`
	Provider      string       `gorm:"type:varchar(255);not null" json:"provider"`
	Configuration util.JSONMap `gorm:"type:jsonb;not null" json:"configuration"`
}

func (LLM) TableName() string {
	return "llms"
}
