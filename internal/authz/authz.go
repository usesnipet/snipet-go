// Package authz holds tenant-membership authorization checks shared by every
// module scoped under /tenants/{tenant_id}/... — factored out once enough
// services needed the same two checks that duplicating them stopped being
// the cheaper option.
package authz

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/model"
)

// RequireMember ensures the caller belongs to the tenant, any role.
func RequireMember(ctx context.Context, tenantID string) (*auth.UserIdentity, error) {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if ok := identity.IsMemberOf(tenantID); ok {
		return identity, nil
	}
	return nil, apperr.Forbidden("not a member of this tenant")
}

// RequireTenantAdmin ensures the caller is an active tenant admin.
func RequireTenantAdmin(ctx context.Context, tenantID string) (*auth.UserIdentity, error) {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if !identity.IsTenantAdmin(tenantID) {
		return nil, apperr.Forbidden("not allowed to manage this tenant")
	}
	return identity, nil
}

// RequireTenantRole ensures the caller has the given tenant role.
func RequireTenantRole(ctx context.Context, tenantID string, role model.MemberRole) (*auth.UserIdentity, error) {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if !identity.IsTenantRole(tenantID, role) {
		return nil, apperr.Forbidden("not allowed to perform this action")
	}
	return identity, nil
}
