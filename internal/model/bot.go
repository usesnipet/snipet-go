package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type BotConfiguration struct {
	LLMs []any `json:"llms"`
}

func (c *BotConfiguration) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan BotConfiguration: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, c)
}

func (c BotConfiguration) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type Bot struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string           `gorm:"type:varchar(255);not null" json:"name"`
	Description   string           `gorm:"type:text;not null" json:"description"`
	Configuration BotConfiguration `gorm:"type:jsonb;not null" json:"configuration"`

	BotMemories []BotMemory `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientBots  []ClientBot `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
