package knowledgeitem

import (
	"time"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type CreateKnowledgeItemDTO struct {
	ExternalID   string       `json:"external_id" validate:"omitempty,max=255"`
	Name         string       `json:"name" validate:"omitempty"`
	Hash         string       `json:"hash" validate:"omitempty,max=128"`
	Metadata     util.JSONMap `json:"metadata"`
	LastModified *time.Time   `json:"last_modified"`
}

type UpdateKnowledgeItemDTO struct {
	ExternalID   *string       `json:"external_id" validate:"omitempty,max=255"`
	Name         *string       `json:"name" validate:"omitempty"`
	Hash         *string       `json:"hash" validate:"omitempty,max=128"`
	Metadata     *util.JSONMap `json:"metadata"`
	LastModified *time.Time    `json:"last_modified"`
}

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
