package auth_provider

import (
	"context"
	"net/http"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
)

type Registry struct {
	providers map[ProviderName]IProvider
}

func NewRegistry() *Registry {
	return &Registry{providers: make(map[ProviderName]IProvider)}
}

func (r *Registry) RegisterProvider(provider IProvider) {
	r.providers[provider.Name()] = provider
}

func (r *Registry) get(ctx context.Context, providerName ProviderName) (IProvider, error) {
	provider, ok := r.providers[providerName]
	if !ok {
		return nil, apperr.NotFound("auth provider not found")
	}
	return provider, nil
}

func (r *Registry) Authenticate(
	ctx context.Context,
	providerName ProviderName,
	appCode string,
	authConfig *model.AppAuthConfig,
	req *http.Request,
) (*Identity, error) {
	provider, err := r.get(ctx, providerName)
	if err != nil {
		return nil, err
	}
	if err := provider.Validate(ctx, appCode, authConfig); err != nil {
		return nil, err
	}
	identity, err := provider.Authenticate(ctx, appCode, authConfig, req)
	if err != nil {
		return nil, err
	}

	if identity.ExternalID == "" {
		return nil, apperr.BadRequest("app authentication response missing external id")
	}

	identity.Provider = providerName
	return identity, nil
}
