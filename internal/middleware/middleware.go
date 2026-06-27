package middleware

import (
	"net/http"
	"time"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/infra/cache"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
)

func APIKeyMiddleware(
	apiKeyService *apikey.Service,
	apiKeyCache cache.ICache,
) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, "API key is required", http.StatusUnauthorized)
				return
			}

			_, found := apiKeyCache.Get(apiKey)
			if found {
				next.ServeHTTP(w, r)
				return
			}
			err := apiKeyService.VerifyAPIKey(r.Context(), apiKey)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, err)
				return
			}
			apiKeyCache.Set(apiKey, true, cache.WithTTL(1*time.Minute))
			next.ServeHTTP(w, r)
		})
	}
}
