package memory

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type FindMemoriesFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindMemoriesFilterDTO) ToFilter() *filter.Options[model.Memory] {
	return filter.New[model.Memory](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type CreateMemoryDTO struct {
	Name          string           `json:"name" validate:"required,max=255"`
	Type          model.MemoryType `json:"type" validate:"required,max=255"`
	IsDefault     bool             `json:"is_default"`
	Provider      string           `json:"provider" validate:"required,max=255"`
	Configuration util.JSONMap     `json:"configuration" validate:"required"`
}

type UpdateMemoryDTO struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
}

type SetAsDefaultMemoryDTO struct {
	MemoryID string `json:"memory_id" validate:"required,uuid"`
}
