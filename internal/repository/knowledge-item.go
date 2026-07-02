package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IKnowledgeItemRepository interface {
	IRepository[model.KnowledgeItem]
}

type KnowledgeItemRepository struct {
	*Repository[model.KnowledgeItem]
}

func NewKnowledgeItemRepository(db *gorm.DB) IKnowledgeItemRepository {
	return &KnowledgeItemRepository{
		Repository: NewRepository[model.KnowledgeItem](db),
	}
}
