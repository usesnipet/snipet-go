package indexedknowledgeitem

import (
	"time"

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
	IndexID         string            `json:"index_id" validate:"required,uuid"`
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
