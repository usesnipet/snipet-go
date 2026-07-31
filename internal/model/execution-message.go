package model

import (
	"github.com/usesnipet/snipet/pkg/msg"
)

type ExecutionMessage struct {
	msg.Message `gorm:"embedded"`

	ExecutionID string    `gorm:"type:uuid;not null;index" json:"execution_id"`
	Execution   Execution `gorm:"foreignKey:ExecutionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
