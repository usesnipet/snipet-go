package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/infra/cache"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
)

func APIKeyMiddleware(apiKeyService *apikey.Service, apiKeyCache cache.ICache) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, "API key is required", http.StatusUnauthorized)
				return
			}

			keyID, found := cache.GetAs[uuid.UUID](apiKeyCache, apiKey)
			if found {
				principal := auth.NewPrincipal(auth.PrincipalTypeAPIKey, &keyID, nil)
				ctx := auth.SetPrincipal(r.Context(), principal)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			key, err := apiKeyService.VerifyAPIKey(r.Context(), apiKey)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, err)
				return
			}
			apiKeyCache.Set(apiKey, key.ID, cache.WithTTL(1*time.Minute))
			principal := auth.NewPrincipal(auth.PrincipalTypeAPIKey, &key.ID, nil)
			ctx := auth.SetPrincipal(r.Context(), principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
