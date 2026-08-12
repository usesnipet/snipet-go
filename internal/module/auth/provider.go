package auth

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
)

// ProviderName identifies a platform-level OAuth provider for tenant-staff
// social login.
type ProviderName string

const (
	ProviderGoogle ProviderName = "google"
	ProviderGithub ProviderName = "github"
)

type Identity struct {
	ExternalID string
	Email      string
	Name       string
	Picture    string
}

type IProvider interface {
	Name() ProviderName
	AuthorizationURL(state string) string
	Exchange(ctx context.Context, code string) (*Identity, error)
}

// ProviderRegistry looks up a configured IProvider by name.
type ProviderRegistry struct {
	providers map[ProviderName]IProvider
}

func NewProviderRegistry(providers ...IProvider) *ProviderRegistry {
	registry := &ProviderRegistry{providers: make(map[ProviderName]IProvider, len(providers))}
	for _, provider := range providers {
		registry.providers[provider.Name()] = provider
	}
	return registry
}

func (r *ProviderRegistry) Get(name ProviderName) (IProvider, error) {
	provider, ok := r.providers[name]
	if !ok {
		return nil, apperr.NotFound("auth provider not found")
	}
	return provider, nil
}
