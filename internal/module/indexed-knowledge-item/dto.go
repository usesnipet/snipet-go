package indexedknowledgeitem

import (
	"time"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type CreateIndexedKnowledgeItemDTO struct {
	Status          model.IndexStatus `json:"status" validate:"omitempty,oneof=pending syncing indexed error"`
	Version         int               `json:"version" validate:"omitempty,min=1"`
	Hash            string            `json:"hash" validate:"omitempty,max=128"`
	IndexedAt       *time.Time        `json:"indexed_at"`
	LastError       string            `json:"last_error"`
	Metadata        util.JSONMap      `json:"metadata"`
	KnowledgeItemID string            `json:"knowledge_item_id" validate:"required,uuid"`
}

type UpdateIndexedKnowledgeItemDTO struct {
	Status    *model.IndexStatus `json:"status" validate:"omitempty,oneof=pending syncing indexed error"`
	Version   *int               `json:"version" validate:"omitempty,min=1"`
	Hash      *string            `json:"hash" validate:"omitempty,max=128"`
	IndexedAt *time.Time         `json:"indexed_at"`
	LastError *string            `json:"last_error"`
	Metadata  *util.JSONMap      `json:"metadata"`
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
