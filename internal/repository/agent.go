package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IAgentRepository interface {
	IRepository[model.Agent]
}

type AgentRepository struct {
	*Repository[model.Agent]
}

func NewAgentRepository(db *gorm.DB) IAgentRepository {
	return &AgentRepository{
		Repository: NewRepository[model.Agent](db),
	}
}
