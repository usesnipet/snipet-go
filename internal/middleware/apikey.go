package middleware

import (
	"net/http"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/infra/cache"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
)

func APIKeyMiddleware(apiKeyService *apikey.Service, apiKeyCache cache.ICache) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKeyAuth(apiKeyService, apiKeyCache, next, w, r)
		})
	}
}
