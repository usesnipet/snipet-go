package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IIndexedKnowledgeItemRepository interface {
	IRepository[model.IndexedKnowledgeItem]
}

type IndexedKnowledgeItemRepository struct {
	*Repository[model.IndexedKnowledgeItem]
}

func NewIndexedKnowledgeItemRepository(db *gorm.DB) IIndexedKnowledgeItemRepository {
	return &IndexedKnowledgeItemRepository{
		Repository: NewRepository[model.IndexedKnowledgeItem](db),
	}
}
