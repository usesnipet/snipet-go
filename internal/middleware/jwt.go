package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
)

func JWT(jwtService *auth.JWTService) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer ") {
				api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")
			claims, err := jwtService.VerifyToken(token)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}

			principal := auth.NewPrincipal(auth.PrincipalTypeJWT, nil, claims)
			ctx := auth.SetPrincipal(r.Context(), principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
