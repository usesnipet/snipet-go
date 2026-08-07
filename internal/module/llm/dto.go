package llm

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// LLMResponse, LLMsPage and DriverInfo exist so swagger annotations in this
// package can reference them without importing internal/model, internal/page
// or pkg/driver directly.
type LLMResponse = model.LLM

type LLMsPage = page.Paginated[model.LLM]

type DriverInfo = driver.Info

type CreateLLMDTO struct {
	Name          string        `json:"name" validate:"required,max=255"`
	Provider      string        `json:"provider" validate:"required,max=255"`
	Configuration jsonx.JSONMap `json:"configuration" validate:"required"`
}

type UpdateLLMDTO struct {
	Name          *string       `json:"name" validate:"omitempty,max=255"`
	Provider      *string       `json:"provider" validate:"omitempty,max=255"`
	Configuration jsonx.JSONMap `json:"configuration" validate:"omitempty"`
}

type FindLLMsFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindLLMsFilterDTO) ToFilter() *filter.Options[model.LLM] {
	return filter.New[model.LLM](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
