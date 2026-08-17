package guard

import (
	"context"
	"net/http"
	"time"

	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/infra/cache"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
)

// cachedAPIKey is what RequireApiKey caches per raw key — just enough to
// rebuild an auth.ApiKeyIdentity on a cache hit without re-verifying.
type cachedAPIKey struct {
	ID       string
	TenantID string
}

// RequireApiKey requires a valid X-API-Key header and sets an
// auth.ApiKeyIdentity.
func RequireApiKey(apiKeyService *apikey.Service, apiKeyCache cache.ICache) Gate {
	return func(r *http.Request) (context.Context, error) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			return nil, auth.ErrNotApplicable
		}

		if cached, found := cache.GetAs[cachedAPIKey](apiKeyCache, key); found {
			return auth.SetApiKeyIdentity(r.Context(), auth.ApiKeyIdentity{APIKeyID: cached.ID, TenantID: cached.TenantID}), nil
		}

		found, err := apiKeyService.VerifyAPIKey(r.Context(), key)
		if err != nil {
			return nil, err
		}
		apiKeyCache.Set(key, cachedAPIKey{ID: found.ID, TenantID: found.TenantID}, cache.WithTTL(1*time.Minute))

		return auth.SetApiKeyIdentity(r.Context(), auth.ApiKeyIdentity{APIKeyID: found.ID, TenantID: found.TenantID}), nil
	}
}
