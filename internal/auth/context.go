package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	apperr "github.com/usesnipet/snipet/internal/app-err"
)

type PrincipalType string

const (
	PrincipalTypeAPIKey      PrincipalType = "api_key"
	PrincipalTypeClientJWT   PrincipalType = "client_jwt"
	PrincipalTypePlatformJWT PrincipalType = "platform_jwt"
)

type Principal struct {
	Type PrincipalType

	APIKeyID  *string
	JWTClaims jwt.Claims
}

func NewPrincipal(t PrincipalType, apiKeyID *string, jwtClaims jwt.Claims) *Principal {
	return &Principal{
		Type:      t,
		APIKeyID:  apiKeyID,
		JWTClaims: jwtClaims,
	}
}

func (p *Principal) GetType() PrincipalType {
	return p.Type
}

func (p *Principal) GetAPIKeyID() *string {
	return p.APIKeyID
}

func (p *Principal) GetJWTClaims() jwt.Claims {
	return p.JWTClaims
}

// Principals is the set of auth methods that succeeded on a request.
type Principals []*Principal

func (ps Principals) OfType(t PrincipalType) (*Principal, bool) {
	for _, p := range ps {
		if p != nil && p.Type == t {
			return p, true
		}
	}
	return nil, false
}

func (ps Principals) Has(t PrincipalType) bool {
	_, ok := ps.OfType(t)
	return ok
}

type principalsKeyType struct{}

var principalsKey = principalsKeyType{}

// SetPrincipals replaces the authenticated-principals list on ctx.
func SetPrincipals(ctx context.Context, principals Principals) context.Context {
	return context.WithValue(ctx, principalsKey, principals)
}

// SetPrincipal is a convenience for tests / single-gate auth: stores a
// one-element principals list.
func SetPrincipal(ctx context.Context, principal *Principal) context.Context {
	return SetPrincipals(ctx, Principals{principal})
}

func GetPrincipals(ctx context.Context) (Principals, bool) {
	principals, ok := ctx.Value(principalsKey).(Principals)
	if !ok || len(principals) == 0 {
		return nil, false
	}
	return principals, true
}

func HasPrincipal(ctx context.Context, t PrincipalType) bool {
	principals, ok := GetPrincipals(ctx)
	return ok && principals.Has(t)
}

func ExpectPrincipal(ctx context.Context, t PrincipalType) (*Principal, error) {
	principals, ok := GetPrincipals(ctx)
	if !ok {
		return nil, apperr.Unauthorized("unauthorized")
	}
	principal, ok := principals.OfType(t)
	if !ok {
		return nil, apperr.Unauthorized("unauthorized")
	}
	return principal, nil
}

func APIKeyID(ctx context.Context) (string, error) {
	principal, err := ExpectPrincipal(ctx, PrincipalTypeAPIKey)
	if err != nil {
		return "", err
	}
	if principal.GetAPIKeyID() == nil || *principal.GetAPIKeyID() == "" {
		return "", apperr.Unauthorized("unauthorized")
	}
	return *principal.GetAPIKeyID(), nil
}

func ClientJWTClaims(ctx context.Context) (*ClientUserClaims, error) {
	principal, err := ExpectPrincipal(ctx, PrincipalTypeClientJWT)
	if err != nil {
		return nil, err
	}
	claims, ok := principal.GetJWTClaims().(*ClientUserClaims)
	if !ok || claims == nil {
		return nil, apperr.Unauthorized("unauthorized")
	}
	return claims, nil
}

func PlatformJWTClaims(ctx context.Context) (*PlatformUserClaims, error) {
	principal, err := ExpectPrincipal(ctx, PrincipalTypePlatformJWT)
	if err != nil {
		return nil, err
	}
	claims, ok := principal.GetJWTClaims().(*PlatformUserClaims)
	if !ok || claims == nil {
		return nil, apperr.Unauthorized("unauthorized")
	}
	return claims, nil
}

func ClientUserID(ctx context.Context) (string, error) {
	claims, err := ClientJWTClaims(ctx)
	if err != nil {
		return "", err
	}
	subject, err := claims.GetSubject()
	if err != nil || subject == "" {
		return "", apperr.Unauthorized("unauthorized")
	}
	return subject, nil
}

func PlatformUserID(ctx context.Context) (string, error) {
	claims, err := PlatformJWTClaims(ctx)
	if err != nil {
		return "", err
	}
	subject, err := claims.GetSubject()
	if err != nil || subject == "" {
		return "", apperr.Unauthorized("unauthorized")
	}
	return subject, nil
}

// ErrNotApplicable means this authenticator's credential was not present
// on the request (e.g. no Bearer header for a JWT gate). Used by
// middleware.Or to skip that gate while still trying the others.
var ErrNotApplicable = errors.New("auth not applicable")
