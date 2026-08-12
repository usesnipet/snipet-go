package middleware

import (
	"net/http"
	"time"

	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/infra/cache"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
)

// RequireAPIKey requires a valid X-API-Key and adds PrincipalTypeAPIKey.
func RequireAPIKey(apiKeyService *apikey.Service, apiKeyCache cache.ICache) AuthMiddleware {
	return single(func(r *http.Request) (*auth.Principal, error) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			return nil, auth.ErrNotApplicable
		}
		keyID, found := cache.GetAs[string](apiKeyCache, apiKey)
		if found {
			return auth.NewPrincipal(auth.PrincipalTypeAPIKey, &keyID, nil), nil
		}
		key, err := apiKeyService.VerifyAPIKey(r.Context(), apiKey)
		if err != nil {
			return nil, err
		}
		apiKeyCache.Set(apiKey, key.ID, cache.WithTTL(1*time.Minute))
		return auth.NewPrincipal(auth.PrincipalTypeAPIKey, &key.ID, nil), nil
	})
}
