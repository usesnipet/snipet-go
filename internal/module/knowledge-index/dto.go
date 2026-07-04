package knowledgeindex

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type CreateKnowledgeIndexDTO struct {
	Name          string       `json:"name" validate:"required,max=255"`
	Driver        string       `json:"driver" validate:"required,max=100"`
	Configuration util.JSONMap `json:"configuration" validate:"required"`
}

type UpdateKnowledgeIndexDTO struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
}

type FilterKnowledgeIndexDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FilterKnowledgeIndexDTO) ToFilter() *filter.Options[model.KnowledgeIndex] {
	return filter.New[model.KnowledgeIndex](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type FilterIndexedKnowledgeItemDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FilterIndexedKnowledgeItemDTO) ToFilter() *filter.Options[model.IndexedKnowledgeItem] {
	return filter.New[model.IndexedKnowledgeItem](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
