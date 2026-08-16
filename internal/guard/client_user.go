package guard

import (
	"context"
	"net/http"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
)

// RequireClientUser requires a valid client-widget end-user JWT
// (Authorization: Bearer) and sets an auth.ClientUserIdentity.
func RequireClientUser(jwtService *auth.JWTService[*auth.ClientUserClaims]) Gate {
	return func(r *http.Request) (context.Context, error) {
		claims, err := verifyBearerJWT(r, jwtService)
		if err != nil {
			return nil, err
		}
		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			return nil, apperr.Unauthorized("unauthorized")
		}

		identity := auth.ClientUserIdentity{UserID: userID, ClientCode: claims.ClientCode}
		return auth.SetClientUserIdentity(r.Context(), identity), nil
	}
}
