package knowledge

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/util"
)

type FilterKnowledgeItemDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FilterKnowledgeItemDTO) ToFilter() *filter.Options[model.KnowledgeItem] {
	return filter.New[model.KnowledgeItem](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type FilterKnowledgeDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FilterKnowledgeDTO) ToFilter() *filter.Options[model.Knowledge] {
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

type TestConnectionDTO struct {
	Driver        string       `json:"driver" validate:"required,max=255"`
	Configuration util.JSONMap `json:"configuration" validate:"required"`
}

type SyncKnowledgeQueryDTO struct {
	Force bool `form:"force" validate:"omitempty"`
}

type DriversDTO struct {
	SourceDrivers []driver.Info `json:"source_drivers"`
}
