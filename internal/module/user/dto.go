package user

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// UserResponse and UsersPage exist so swagger annotations in this package can
// reference them without importing internal/model or internal/page directly.
type UserResponse = model.User

type UsersPage = page.Paginated[model.User]

type CreateAnonymousClientUserDTO struct {
	Name     *string       `json:"name" validate:"omitempty,max=255"`
	Metadata jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type CreateAuthenticatedClientUserDTO struct {
	ExternalID string        `json:"external_id" validate:"required,max=255"`
	Name       string        `json:"name" validate:"omitempty,max=255"`
	Email      string        `json:"email" validate:"required,email"`
	Picture    *string       `json:"picture" validate:"omitempty,url"`
	Metadata   jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type FindUsersFilterDTO struct {
	NameOrder *filter.OrderDirection `form:"name_order" validate:"omitempty,oneof=asc desc"`
	Take      *int                   `form:"take" validate:"omitempty,min=1"`
	Skip      *int                   `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindUsersFilterDTO) ToFilter() *filter.Options[model.User] {
	return filter.New[model.User](
		filter.PtrOrderBy("name", dto.NameOrder),
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
