package middleware

import (
	"net/http"
	"strings"

	"github.com/usesnipet/snipet/internal/auth"
)

// RequireClientJWT requires a valid client-widget JWT
// (Authorization: Bearer) and adds PrincipalTypeClientJWT.
func RequireClientJWT(jwtService *auth.JWTService[*auth.ClientUserClaims]) AuthMiddleware {
	return single(func(r *http.Request) (*auth.Principal, error) {
		return authenticateJWT(r, jwtService, auth.PrincipalTypeClientJWT)
	})
}

// RequirePlatformJWT requires a valid tenant-staff JWT
// (Authorization: Bearer) and adds PrincipalTypePlatformJWT.
func RequirePlatformJWT(jwtService *auth.JWTService[*auth.PlatformUserClaims]) AuthMiddleware {
	return single(func(r *http.Request) (*auth.Principal, error) {
		return authenticateJWT(r, jwtService, auth.PrincipalTypePlatformJWT)
	})
}

func authenticateJWT[T auth.Claims](
	r *http.Request,
	jwtService *auth.JWTService[T],
	principalType auth.PrincipalType,
) (*auth.Principal, error) {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return nil, auth.ErrNotApplicable
	}

	claims, err := jwtService.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	return auth.NewPrincipal(principalType, nil, claims), nil
}
