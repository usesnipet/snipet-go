package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IKnowledgeRepository interface {
	IRepository[model.Knowledge]
}

type KnowledgeRepository struct {
	*Repository[model.Knowledge]
}

func NewKnowledgeRepository(db *gorm.DB) IKnowledgeRepository {
	return &KnowledgeRepository{
		Repository: NewRepository[model.Knowledge](db),
	}
}
