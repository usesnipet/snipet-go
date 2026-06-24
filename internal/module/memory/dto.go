package memory

import "encoding/json"

type CreateMemoryDTO struct {
	Name          string          `json:"name" validate:"required,max=255"`
	Type          string          `json:"type" validate:"required,max=255"`
	IsDefault     bool            `json:"is_default"`
	Provider      string          `json:"provider" validate:"required,max=255"`
	Configuration json.RawMessage `json:"configuration" validate:"required"`
}

type UpdateMemoryDTO struct {
	Name          *string          `json:"name" validate:"omitempty,max=255"`
	Type          *string          `json:"type" validate:"omitempty,max=255"`
	IsDefault     *bool            `json:"is_default"`
	Provider      *string          `json:"provider" validate:"omitempty,max=255"`
	Configuration *json.RawMessage `json:"configuration"`
}
