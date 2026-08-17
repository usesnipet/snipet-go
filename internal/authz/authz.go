// Package authz holds tenant-membership authorization checks shared by every
// module scoped under /tenants/{tenant_id}/... — factored out once enough
// services needed the same two checks that duplicating them stopped being
// the cheaper option.
package authz

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
)

// RequireMember ensures the caller belongs to the tenant, any role.
func RequireMember(ctx context.Context, tenantID string) (*auth.UserIdentity, error) {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if !identity.IsMemberOf(tenantID) {
		return nil, apperr.Forbidden("not a member of this tenant")
	}
	return identity, nil
}

// RequireAdmin ensures the caller is an active tenant admin.
func RequireAdmin(ctx context.Context, tenantID string) (*auth.UserIdentity, error) {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if !identity.IsTenantAdmin(tenantID) {
		return nil, apperr.Forbidden("not allowed to manage this tenant")
	}
	return identity, nil
}
