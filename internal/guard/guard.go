// Package guard authenticates requests and stashes the resulting identity
// (see internal/auth's UserIdentity, ApiKeyIdentity, ClientUserIdentity) on
// the request context. Each Require* constructor below is one auth scheme;
// wire the one an endpoint needs via .Handler(), or accept several with Or.
package guard

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
)

// Gate is a single authentication attempt against a request. On success it
// returns a context carrying whatever identity it authenticated (via
// auth.SetUserIdentity / SetApiKeyIdentity / SetClientUserIdentity). On
// failure it returns auth.ErrNotApplicable when its credential was simply
// absent from the request — so Or can fall through to the next gate — or
// any other error when the credential was present but invalid.
type Gate func(r *http.Request) (context.Context, error)

// Handler turns a Gate into chi-compatible middleware that requires it to
// succeed.
func (g Gate) Handler() api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, err := g(r)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Or tries each gate in order and keeps the context from the first that
// succeeds — every successful gate sets its own identity, so which one fired
// is recoverable downstream (auth.CurrentApiKey / CurrentClientUser / ...).
// Gates reporting auth.ErrNotApplicable (credential absent) are skipped; a
// gate whose credential is present but invalid fails the whole chain
// immediately. At least one gate must succeed.
func Or(gates ...Gate) Gate {
	return func(r *http.Request) (context.Context, error) {
		for _, g := range gates {
			ctx, err := g(r)
			if errors.Is(err, auth.ErrNotApplicable) {
				continue
			}
			if err != nil {
				return nil, err
			}
			return ctx, nil
		}
		return nil, errors.New("unauthorized")
	}
}

// verifyBearerJWT extracts and verifies an "Authorization: Bearer ..."
// token. Returns auth.ErrNotApplicable when the header is absent, so gates
// built on it compose cleanly with Or.
func verifyBearerJWT[T auth.Claims](r *http.Request, jwtService *auth.JWTService[T]) (T, error) {
	var zero T

	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return zero, auth.ErrNotApplicable
	}

	return jwtService.VerifyToken(token)
}
