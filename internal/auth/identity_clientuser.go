package auth

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
)

// ClientUserIdentity is the authenticated client-widget end-user (bearer
// JWT), set by guard.RequireClientUser.
type ClientUserIdentity struct {
	UserID     string
	ClientCode string
}

type clientUserIdentityKeyType struct{}

var clientUserIdentityKey = clientUserIdentityKeyType{}

// SetClientUserIdentity stores the authenticated ClientUserIdentity on ctx.
func SetClientUserIdentity(ctx context.Context, identity ClientUserIdentity) context.Context {
	return context.WithValue(ctx, clientUserIdentityKey, identity)
}

// CurrentClientUser returns the ClientUserIdentity loaded by
// guard.RequireClientUser for this request.
func CurrentClientUser(ctx context.Context) (ClientUserIdentity, error) {
	identity, ok := ctx.Value(clientUserIdentityKey).(ClientUserIdentity)
	if !ok || identity.UserID == "" {
		return ClientUserIdentity{}, apperr.Unauthorized("unauthorized")
	}
	return identity, nil
}

// HasClientUser reports whether the request authenticated as a client
// end-user.
func HasClientUser(ctx context.Context) bool {
	_, err := CurrentClientUser(ctx)
	return err == nil
}
