package auth_provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

const OIDCProviderName ProviderName = "oidc"

type OIDCProvider struct{}

type oidcClaims struct {
	jwt.Claims
	Issuer string `json:"iss"`

	Subject string  `json:"sub"`
	Name    *string `json:"name"`
}

func NewOIDCProvider() IProvider {
	return &OIDCProvider{}
}

func (p *OIDCProvider) Name() ProviderName {
	return OIDCProviderName
}

func (p *OIDCProvider) Validate(
	ctx context.Context,
	appCode string,
	authConfig *model.AppAuthConfig,
) error {
	if !authConfig.OIDC.Enabled {
		return apperr.BadRequest("app oidc is not enabled")
	}

	return nil
}

func (p *OIDCProvider) Authenticate(
	ctx context.Context,
	appCode string,
	authConfig *model.AppAuthConfig,
	req *http.Request,
) (*Identity, error) {

	auth := req.Header.Get("Authorization")
	if auth == "" {
		return nil, apperr.Unauthorized("missing authorization header")
	}

	const prefix = "Bearer "

	if !strings.HasPrefix(auth, prefix) {
		return nil, apperr.Unauthorized("invalid authorization header")
	}

	token := strings.TrimPrefix(auth, prefix)

	// Lê as claims sem validar para descobrir o issuer
	var unverified oidcClaims

	_, _, err := new(jwt.Parser).ParseUnverified(token, &unverified)
	if err != nil {
		return nil, apperr.BadRequest("invalid oidc token")
	}

	issuer := unverified.Issuer
	if authConfig.OIDC.Issuer != "" {
		issuer = authConfig.OIDC.Issuer
	}

	if issuer == "" {
		return nil, apperr.BadRequest("missing oidc issuer")
	}

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, apperr.BadRequest("invalid oidc issuer")
	}

	config := &oidc.Config{}

	if authConfig.OIDC.Audience != "" {
		config.ClientID = authConfig.OIDC.Audience
	} else {
		config.SkipClientIDCheck = true
	}

	verifier := provider.Verifier(config)

	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, apperr.Unauthorized("invalid oidc token")
	}

	var claims oidcClaims

	if err := idToken.Claims(&claims); err != nil {
		return nil, apperr.BadRequest("invalid oidc claims")
	}

	metadata, err := jsonx.ToJSONMap(claims)
	if err != nil {
		return nil, apperr.BadRequest("provider metadata conversion failed")
	}

	return &Identity{
		Provider:   OIDCProviderName,
		ExternalID: claims.Subject,
		Name:       claims.Name,
		Metadata:   metadata,
	}, nil
}
