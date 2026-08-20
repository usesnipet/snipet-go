package auth

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
)

// AppKeyIdentity is the App that authenticated the request via its own key,
// set by guard.RequireAppKey.
type AppKeyIdentity struct {
	AppID string
	Code  string
}

func (i AppKeyIdentity) Is(codeOrId string) error {
	if codeOrId == i.AppID || codeOrId == i.Code {
		return nil
	}
	return apperr.Forbidden("this app key does not have access to this app")
}

type appKeyIdentityKeyType struct{}

var appKeyIdentityKey = appKeyIdentityKeyType{}

// SetAppKeyIdentity stores the authenticated AppKeyIdentity on ctx.
func SetAppKeyIdentity(ctx context.Context, identity AppKeyIdentity) context.Context {
	return context.WithValue(ctx, appKeyIdentityKey, identity)
}

// CurrentAppKey returns the AppKeyIdentity loaded by guard.RequireAppKey for
// this request.
func CurrentAppKey(ctx context.Context) (AppKeyIdentity, error) {
	identity, ok := ctx.Value(appKeyIdentityKey).(AppKeyIdentity)
	if !ok || identity.AppID == "" {
		return AppKeyIdentity{}, apperr.Unauthorized("unauthorized")
	}
	return identity, nil
}

// HasAppKey reports whether the request authenticated via an App key.
func HasAppKey(ctx context.Context) bool {
	_, err := CurrentAppKey(ctx)
	return err == nil
}
