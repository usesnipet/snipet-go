package memory

import (
	"encoding/json"

	"github.com/usesnipet/snipet/internal/model"
)

type CreateMemoryDTO struct {
	Name          string           `json:"name" validate:"required,max=255"`
	Type          model.MemoryType `json:"type" validate:"required,max=255"`
	IsDefault     bool             `json:"is_default" validate:"required"`
	Provider      string           `json:"provider" validate:"required,max=255"`
	Configuration json.RawMessage  `json:"configuration" validate:"required"`
}

type UpdateMemoryDTO struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
}

type SetAsDefaultMemoryDTO struct {
	MemoryID string `json:"memory_id" validate:"required,uuid"`
}
