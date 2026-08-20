package auth

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
)

// AppUserIdentity is the authenticated app end-user (bearer
// JWT), set by guard.RequireAppUser.
type AppUserIdentity struct {
	UserID  string
	AppCode string
}

func (i AppUserIdentity) CanAccessApp(appCode string) error {
	if i.AppCode == appCode {
		return nil
	}
	return apperr.Forbidden("app code mismatch")
}

type appUserIdentityKeyType struct{}

var appUserIdentityKey = appUserIdentityKeyType{}

// SetAppUserIdentity stores the authenticated AppUserIdentity on ctx.
func SetAppUserIdentity(ctx context.Context, identity AppUserIdentity) context.Context {
	return context.WithValue(ctx, appUserIdentityKey, identity)
}

// CurrentAppUser returns the AppUserIdentity loaded by
// guard.RequireAppUser for this request.
func CurrentAppUser(ctx context.Context) (AppUserIdentity, error) {
	identity, ok := ctx.Value(appUserIdentityKey).(AppUserIdentity)
	if !ok || identity.UserID == "" {
		return AppUserIdentity{}, apperr.Unauthorized("unauthorized")
	}
	return identity, nil
}

// HasAppUser reports whether the request authenticated as a app
// end-user.
func HasAppUser(ctx context.Context) bool {
	_, err := CurrentAppUser(ctx)
	return err == nil
}
