package auth

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
)

// UserIdentity is the authenticated platform user's profile plus every
// tenant membership they hold. guard.RequireUser loads it once per request
// (verifying the bearer JWT and resolving the user + memberships in one
// pass) and stashes it on the context, so services read it via
// auth.CurrentUser instead of each re-querying the user/member
// repositories.
type UserIdentity struct {
	User        *model.User
	Memberships []*model.Member
}

// MembershipOf returns the caller's membership row for tenantID, if any.
func (i *UserIdentity) MembershipOf(tenantID string) (*model.Member, bool) {
	for _, m := range i.Memberships {
		if m != nil && m.TenantID == tenantID {
			return m, true
		}
	}
	return nil, false
}

// IsMemberOf reports whether the caller belongs to tenantID, in any role.
func (i *UserIdentity) IsMemberOf(tenantID string) bool {
	_, ok := i.MembershipOf(tenantID)
	return ok
}

// IsTenantAdmin reports whether the caller is an active admin member of
// tenantID. It does not consider platform-admin status — combine with
// IsPlatformAdmin where a platform-admin override is wanted.
func (i *UserIdentity) IsTenantAdmin(tenantID string) bool {
	member, ok := i.MembershipOf(tenantID)
	return ok && member.IsActive && member.Role == model.RoleAdmin
}

func (i *UserIdentity) IsTenantRole(tenantID string, role model.Role) bool {
	member, ok := i.MembershipOf(tenantID)
	return ok && member.Role == role
}

// IsPlatformAdmin reports whether the caller is a platform-wide admin
// (model.User.IsAdmin), independent of tenant membership.
func (i *UserIdentity) IsPlatformAdmin() bool {
	return i.User != nil && i.User.IsAdmin
}

type userIdentityKeyType struct{}

var userIdentityKey = userIdentityKeyType{}

// SetUserIdentity stores the loaded UserIdentity on ctx.
func SetUserIdentity(ctx context.Context, identity *UserIdentity) context.Context {
	return context.WithValue(ctx, userIdentityKey, identity)
}

// CurrentUser returns the UserIdentity loaded by guard.RequireUser for this
// request.
func CurrentUser(ctx context.Context) (*UserIdentity, error) {
	identity, ok := ctx.Value(userIdentityKey).(*UserIdentity)
	if !ok || identity == nil {
		return nil, apperr.Unauthorized("unauthorized")
	}
	return identity, nil
}
