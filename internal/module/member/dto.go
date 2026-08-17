package member

import (
	"time"

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

// CreateMemberDTO is only accepted on unlicensed (single-tenant) instances —
// see member.Service.Create — since there's no email-based invitation flow
// to reach a Tenant through without one.
type CreateMemberDTO struct {
	Name            string     `json:"name" validate:"required,max=255"`
	Email           string     `json:"email" validate:"required,email,max=255"`
	Password        string     `json:"password" validate:"required,min=8"`
	ConfirmPassword string     `json:"confirm_password" validate:"required,eqfield=Password"`
	Role            model.Role `json:"role" validate:"required,oneof=admin user"`
}

// invitationStatusExpired is a filter-only pseudo-status: invitations that
// are still pending but whose expiration date has passed. It has no
// corresponding model.InvitationStatus, since expiry is derived from
// ExpiresAt rather than stored.
const invitationStatusExpired = "expired"

type FindInvitationsFilterDTO struct {
	Status *string `form:"status" validate:"omitempty,oneof=pending accepted declined expired"`
	Take   *int    `form:"take" validate:"omitempty,min=1"`
	Skip   *int    `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindInvitationsFilterDTO) ToFilter() *filter.Options[model.TenantInvitation] {
	opts := []filter.Option{
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	}
	if dto.Status != nil {
		switch *dto.Status {
		case invitationStatusExpired:
			opts = append(opts,
				filter.WhereEq("status", model.InvitationStatusPending),
				filter.WhereLt("expires_at", time.Now()),
			)
		case string(model.InvitationStatusPending):
			// "Pending" excludes invitations that are pending but past their
			// expiration date; those surface under the "expired" filter instead.
			opts = append(opts,
				filter.WhereEq("status", model.InvitationStatusPending),
				filter.WhereGte("expires_at", time.Now()),
			)
		default:
			opts = append(opts, filter.WhereEq("status", model.InvitationStatus(*dto.Status)))
		}
	}
	return filter.New[model.TenantInvitation](opts...)
}

// InvitationInfoResponse describes an invitation looked up by its token,
// together with the tenant it belongs to. Used by the public invite-preview
// page so a user can see what they're accepting before authenticating.
type InvitationInfoResponse struct {
	Invite *model.TenantInvitation `json:"invite"`
	Tenant *model.Tenant           `json:"tenant"`
}

type AcceptInvitationDTO struct {
	Token string `json:"token" validate:"required"`
}

type DeclineInvitationDTO struct {
	Token string `json:"token" validate:"required"`
}
