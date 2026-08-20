package app

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
)

// AppResponse and AppsPage exist so swagger annotations in this package can
// reference them without importing internal/model or internal/page directly.
type AppResponse = model.App

type AppsPage = page.Paginated[model.App]

type CreateAppDTO struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateAppDTO struct {
	Name        *string `json:"name" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

type UpdateAppAuthConfigDTO struct {
	model.AppAuthConfig
}

func (dto *UpdateAppAuthConfigDTO) ToModel() model.AppAuthConfig {
	return dto.AppAuthConfig
}

type AppWithSecret struct {
	*model.App
	Key string `json:"key"`
}

type FindAppsFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindAppsFilterDTO) ToFilter() *filter.Options[model.App] {
	return filter.New[model.App](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
