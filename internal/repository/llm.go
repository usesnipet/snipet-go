package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type ILLMRepository interface {
	IRepository[model.LLM]
}

type LLMRepository struct {
	*Repository[model.LLM]
}

func NewLLMRepository(db *gorm.DB) ILLMRepository {
	return &LLMRepository{
		Repository: NewRepository[model.LLM](db),
	}
}
