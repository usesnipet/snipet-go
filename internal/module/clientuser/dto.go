package clientuser

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// ClientUserResponse and UsersPage exist so swagger annotations in this package can
// reference them without importing internal/model or internal/page directly.
type ClientUserResponse = model.ClientUser

type ClientUsersPage = page.Paginated[model.ClientUser]

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

type FindClientUsersFilterDTO struct {
	NameOrder *filter.OrderDirection `form:"name_order" validate:"omitempty,oneof=asc desc"`
	Take      *int                   `form:"take" validate:"omitempty,min=1"`
	Skip      *int                   `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindClientUsersFilterDTO) ToFilter() *filter.Options[model.ClientUser] {
	return filter.New[model.ClientUser](
		filter.PtrOrderBy("name", dto.NameOrder),
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
