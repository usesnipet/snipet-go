package middleware

import (
	"errors"
	"net/http"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
)

// AuthMiddleware is a composable auth gate. Use Handler() with chi's
// r.Use, or pass into Or to accept any combination of schemes — every
// gate that successfully authenticates contributes a Principal to the
// request's principals list.
type AuthMiddleware struct {
	authenticate func(r *http.Request) (auth.Principals, error)
}

// Handler returns a chi middleware that requires authentication and
// stores the resulting principals list on the request context.
func (m AuthMiddleware) Handler() api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principals, err := m.authenticate(r)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}
			ctx := auth.SetPrincipals(r.Context(), principals)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func single(try func(r *http.Request) (*auth.Principal, error)) AuthMiddleware {
	return AuthMiddleware{
		authenticate: func(r *http.Request) (auth.Principals, error) {
			principal, err := try(r)
			if err != nil {
				return nil, err
			}
			return auth.Principals{principal}, nil
		},
	}
}

// Or runs every gate. Each successful authentication is appended to the
// principals list. Gates whose credential is absent (ErrNotApplicable)
// are skipped. A gate whose credential is present but invalid rejects
// the whole request. At least one success is required.
func Or(gates ...AuthMiddleware) AuthMiddleware {
	return AuthMiddleware{
		authenticate: func(r *http.Request) (auth.Principals, error) {
			var principals auth.Principals
			for _, gate := range gates {
				ps, err := gate.authenticate(r)
				if errors.Is(err, auth.ErrNotApplicable) {
					continue
				}
				if err != nil {
					return nil, err
				}
				principals = append(principals, ps...)
			}
			if len(principals) == 0 {
				return nil, errors.New("unauthorized")
			}
			return principals, nil
		},
	}
}
