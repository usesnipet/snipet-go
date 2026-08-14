package auth

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
)

// Identity is the authenticated platform user's profile plus every tenant
// membership they hold. middleware.RequireIdentity loads it once per
// request and stashes it on the context, so services read it via
// CurrentIdentity instead of each re-querying the user/member repositories.
type Identity struct {
	User        *model.User
	Memberships []*model.Member
}

// MembershipOf returns the caller's membership row for tenantID, if any.
func (i *Identity) MembershipOf(tenantID string) (*model.Member, bool) {
	for _, m := range i.Memberships {
		if m != nil && m.TenantID == tenantID {
			return m, true
		}
	}
	return nil, false
}

// IsMemberOf reports whether the caller belongs to tenantID, in any role.
func (i *Identity) IsMemberOf(tenantID string) bool {
	_, ok := i.MembershipOf(tenantID)
	return ok
}

// IsTenantAdmin reports whether the caller is an active admin member of
// tenantID. It does not consider platform-admin status — combine with
// IsPlatformAdmin where a platform-admin override is wanted.
func (i *Identity) IsTenantAdmin(tenantID string) bool {
	member, ok := i.MembershipOf(tenantID)
	return ok && member.IsActive && member.Role == model.RoleAdmin
}

// IsPlatformAdmin reports whether the caller is a platform-wide admin
// (model.User.IsAdmin), independent of tenant membership.
func (i *Identity) IsPlatformAdmin() bool {
	return i.User != nil && i.User.IsAdmin
}

type identityKeyType struct{}

var identityKey = identityKeyType{}

// SetIdentity stores the loaded Identity on ctx.
func SetIdentity(ctx context.Context, identity *Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

// CurrentIdentity returns the Identity loaded by middleware.RequireIdentity
// for this request.
func CurrentIdentity(ctx context.Context) (*Identity, error) {
	identity, ok := ctx.Value(identityKey).(*Identity)
	if !ok || identity == nil {
		return nil, apperr.Unauthorized("unauthorized")
	}
	return identity, nil
}

// CurrentUser is a shorthand for CurrentIdentity(ctx).User.
func CurrentUser(ctx context.Context) (*model.User, error) {
	identity, err := CurrentIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return identity.User, nil
}
