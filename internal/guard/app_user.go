package guard

import (
	"context"
	"net/http"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
)

// RequireAppUser requires a valid app end-user JWT
// (Authorization: Bearer) and sets an auth.AppUserIdentity.
func RequireAppUser(jwtService *auth.JWTService[*auth.AppUserClaims]) Gate {
	return func(r *http.Request) (context.Context, error) {
		claims, err := verifyBearerJWT(r, jwtService)
		if err != nil {
			return nil, err
		}
		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			return nil, apperr.Unauthorized("unauthorized")
		}

		identity := auth.AppUserIdentity{UserID: userID, AppCode: claims.AppCode}
		return auth.SetAppUserIdentity(r.Context(), identity), nil
	}
}
