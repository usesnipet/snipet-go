package knowledge

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type KnowledgeFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *KnowledgeFilterDTO) ToFilter() *filter.Options[model.Knowledge] {
	return filter.New[model.Knowledge](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type CreateKnowledgeDTO struct {
	Name          string       `json:"name" validate:"required,max=255"`
	Description   string       `json:"description" validate:"omitempty"`
	Driver        string       `json:"driver" validate:"required,max=255"`
	Configuration util.JSONMap `json:"configuration" validate:"required"`
}

type UpdateKnowledgeDTO struct {
	Name        *string `json:"name" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty"`
}
