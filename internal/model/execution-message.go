package model

import (
	"github.com/usesnipet/snipet/pkg/msg"
)

type ExecutionMessage struct {
	msg.Message `gorm:"embedded"`

	TenantID    string     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ExecutionID string     `gorm:"type:uuid;not null;index" json:"execution_id"`
	Tenant      *Tenant    `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"tenant"`
	Execution   *Execution `gorm:"foreignKey:ExecutionID;references:ID;constraint:OnDelete:CASCADE" json:"execution"`
}
