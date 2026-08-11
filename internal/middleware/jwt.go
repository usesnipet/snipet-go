package middleware

import (
	"net/http"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	clientauth "github.com/usesnipet/snipet/internal/module/auth"
)

func JWT(jwtService *auth.JWTService[*clientauth.UserClaims]) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtAuth(jwtService, next, w, r)
		})
	}
}
