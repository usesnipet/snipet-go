package model

import "github.com/usesnipet/snipet/pkg/jsonx"

type LLM struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	TenantID      string        `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name          string        `gorm:"type:varchar(255);not null" json:"name"`
	Provider      string        `gorm:"type:varchar(255);not null" json:"provider"`
	Configuration jsonx.JSONMap `gorm:"type:jsonb;not null" json:"configuration"`

	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (LLM) TableName() string {
	return "llms"
}
