package repository

import (
	"context"
	"errors"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IAgentRepository interface {
	IRepository[model.Agent]
	ReplaceLLMs(ctx context.Context, agentID string, llmIDs []string) error
}

type AgentRepository struct {
	*Repository[model.Agent]
}

func NewAgentRepository(db *gorm.DB) IAgentRepository {
	return &AgentRepository{
		Repository: NewRepository[model.Agent](db),
	}
}

func (r *AgentRepository) withLLMs(db *gorm.DB) *gorm.DB {
	return db.
		Preload("AgentToLLMs", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority ASC")
		}).
		Preload("AgentToLLMs.LLM")
}

func (r *AgentRepository) Filter(
	ctx context.Context,
	filterOptions *filter.Options[model.Agent],
) (*page.Paginated[model.Agent], error) {
	if filterOptions == nil {
		filterOptions = filter.Default[model.Agent]()
	}

	db := r.db(ctx)
	var total int64
	if err := db.Model(&model.Agent{}).Count(&total).Error; err != nil {
		return nil, err
	}

	chain, err := filterOptions.ToGormTx(db.Model(&model.Agent{}))
	if err != nil {
		return nil, err
	}

	var data []model.Agent
	if err := r.withLLMs(chain).Find(&data).Error; err != nil {
		return nil, err
	}
	return page.NewPaginated(data, total, int64(filterOptions.Skip), int64(filterOptions.Take)), nil
}

func (r *AgentRepository) FindByID(ctx context.Context, id string) (*model.Agent, error) {
	agents, err := r.Filter(ctx, filter.New[model.Agent](
		filter.WhereEq("id", id),
	))

	if errors.Is(err, gorm.ErrRecordNotFound) || agents.IsEmpty() {
		return nil, apperr.NotFound("entity not found")
	}
	if err != nil {
		return nil, err
	}
	return agents.First(), nil
}

func (r *AgentRepository) ReplaceLLMs(ctx context.Context, agentID string, llmIDs []string) error {
	db := r.db(ctx)
	if err := db.Where("agent_id = ?", agentID).Delete(&model.AgentToLLM{}).Error; err != nil {
		return err
	}
	if len(llmIDs) == 0 {
		return nil
	}

	rels := make([]model.AgentToLLM, len(llmIDs))
	for i, llmID := range llmIDs {
		rels[i] = model.AgentToLLM{
			AgentID:  agentID,
			LLMID:    llmID,
			Priority: i,
		}
	}
	return db.Create(&rels).Error
}
