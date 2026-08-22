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
	Public      bool   `json:"public"`
}

type UpdateAppDTO struct {
	Name        *string `json:"name" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Public      *bool   `json:"public"`
}

type UpdateAppAuthConfigDTO struct {
	AuthConfig model.AppAuthConfig `json:"auth_config"`
}

func (dto *UpdateAppAuthConfigDTO) ToModel() model.AppAuthConfig {
	return dto.AuthConfig
}

type AppWithSecret struct {
	*model.App
	Key string `json:"key"`
}

// PublicAppDTO is what an unauthenticated caller (e.g. a frontend-only app's
// widget) may learn about an app — deliberately excludes anything internal
// like status, tenant, or key material.
type PublicAppDTO struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
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
