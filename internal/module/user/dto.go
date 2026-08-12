package user

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
)

// UserResponse and UsersPage exist so swagger annotations in this package
// can reference them without importing internal/model or internal/page directly.
type UserResponse = model.User

type UsersPage = page.Paginated[model.User]

type CreateUserDTO struct {
	Name     string  `json:"name" validate:"required,max=255"`
	Email    string  `json:"email" validate:"required,email,max=255"`
	Password string  `json:"password" validate:"required,min=8"`
	Picture  *string `json:"picture" validate:"omitempty,max=255"`
	IsAdmin  bool    `json:"is_admin"`
}

type UpdateUserDTO struct {
	Name    *string `json:"name" validate:"omitempty,max=255"`
	Email   *string `json:"email" validate:"omitempty,email,max=255"`
	Picture *string `json:"picture" validate:"omitempty,max=255"`
	IsAdmin *bool   `json:"is_admin" validate:"omitempty"`
}

type FindUsersFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindUsersFilterDTO) ToFilter() *filter.Options[model.User] {
	return filter.New[model.User](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type UpdateProfilePictureDTO struct {
	Picture string `json:"picture" validate:"required,max=255"`
}
