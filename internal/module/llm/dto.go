package llm

import (
	"github.com/usesnipet/snipet/internal/util"
)

type CreateLLMDTO struct {
	Name          string       `json:"name" validate:"required,max=255"`
	Provider      string       `json:"provider" validate:"required,max=255"`
	Configuration util.JSONMap `json:"configuration" validate:"required"`
}

type UpdateLLMDTO struct {
	Name          *string      `json:"name" validate:"omitempty,max=255"`
	Provider      *string      `json:"provider" validate:"omitempty,max=255"`
	Configuration util.JSONMap `json:"configuration" validate:"omitempty"`
}
