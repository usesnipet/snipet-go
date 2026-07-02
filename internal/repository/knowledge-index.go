package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IKnowledgeIndexRepository interface {
	IRepository[model.KnowledgeIndex]
}

type KnowledgeIndexRepository struct {
	*Repository[model.KnowledgeIndex]
}

func NewKnowledgeIndexRepository(db *gorm.DB) IKnowledgeIndexRepository {
	return &KnowledgeIndexRepository{
		Repository: NewRepository[model.KnowledgeIndex](db),
	}
}
