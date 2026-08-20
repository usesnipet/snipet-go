package appuser

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// AppUserResponse and AppUsersPage exist so swagger annotations in this
// package can reference them without importing internal/model or
// internal/page directly.
type AppUserResponse = model.AppUser

type AppUsersPage = page.Paginated[model.AppUser]

type CreateAnonymousAppUserDTO struct {
	Name     *string       `json:"name" validate:"omitempty,max=255"`
	Metadata jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type CreateAppUserDTO struct {
	ExternalID string        `json:"external_id" validate:"required,max=255"`
	Name       string        `json:"name" validate:"omitempty,max=255"`
	Email      string        `json:"email" validate:"required,email"`
	Picture    *string       `json:"picture" validate:"omitempty,url"`
	Metadata   jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type FindAppUsersFilterDTO struct {
	NameOrder *filter.OrderDirection `form:"name_order" validate:"omitempty,oneof=asc desc"`
	Take      *int                   `form:"take" validate:"omitempty,min=1"`
	Skip      *int                   `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindAppUsersFilterDTO) ToFilter() *filter.Options[model.AppUser] {
	return filter.New[model.AppUser](
		filter.PtrOrderBy("name", dto.NameOrder),
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
