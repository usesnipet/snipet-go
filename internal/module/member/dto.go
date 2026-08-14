package member

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
)

// InvitationResponse and InvitationsPage exist so swagger annotations in this
// package can reference them without importing internal/model or
// internal/page directly.
type InvitationResponse = model.TenantInvitation

type InvitationsPage = page.Paginated[model.TenantInvitation]

type MemberResponse = model.Member

type MembersPage page.Paginated[model.Member]

type FindMembersFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindMembersFilterDTO) ToFilter() *filter.Options[model.Member] {
	return filter.New[model.Member](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type UpdateMemberRoleDTO struct {
	Role model.Role `json:"role" validate:"required,oneof=admin user"`
}

type InviteMemberDTO struct {
	Email string     `json:"email" validate:"required,email,max=255"`
	Role  model.Role `json:"role" validate:"required,oneof=admin user"`
}

type FindInvitationsFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindInvitationsFilterDTO) ToFilter() *filter.Options[model.TenantInvitation] {
	return filter.New[model.TenantInvitation](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type AcceptInvitationDTO struct {
	Token string `json:"token" validate:"required"`
}
