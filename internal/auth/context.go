package auth

import (
	"context"

	"github.com/google/uuid"
)

type PrincipalType string

const (
	PrincipalTypeAPIKey PrincipalType = "api_key"
	PrincipalTypeJWT    PrincipalType = "jwt"
)

type Principal struct {
	Type PrincipalType

	APIKeyID  *uuid.UUID
	JWTClaims *UserClaims
}

func NewPrincipal(t PrincipalType, apiKeyID *uuid.UUID, jwtClaims *UserClaims) *Principal {
	return &Principal{
		Type:      t,
		APIKeyID:  apiKeyID,
		JWTClaims: jwtClaims,
	}
}

func (p *Principal) GetType() PrincipalType {
	return p.Type
}

func (p *Principal) GetAPIKeyID() *uuid.UUID {
	return p.APIKeyID
}

func (p *Principal) GetJWTClaims() *UserClaims {
	return p.JWTClaims
}

const principalKey = "principal"

func SetPrincipal(ctx context.Context, principal *Principal) context.Context {
	return context.WithValue(ctx, principalKey, principal)
}

func GetPrincipal(ctx context.Context) (*Principal, bool) {
	principal, ok := ctx.Value(principalKey).(*Principal)
	if !ok {
		return nil, false
	}
	return principal, true
}
