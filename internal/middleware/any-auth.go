package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/infra/cache"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
)

func jwtAuth(jwtService *auth.JWTService, next http.Handler, w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	claims, err := jwtService.VerifyToken(token)
	if err != nil {
		api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	principal := auth.NewPrincipal(auth.PrincipalTypeJWT, nil, claims)
	ctx := auth.SetPrincipal(r.Context(), principal)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func apiKeyAuth(apiKeyService *apikey.Service, apiKeyCache cache.ICache, next http.Handler, w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		http.Error(w, "API key is required", http.StatusUnauthorized)
		return
	}
	keyID, found := cache.GetAs[string](apiKeyCache, apiKey)
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
}

func AnyAuth(jwtService *auth.JWTService, apiKeyService *apikey.Service, apiKeyCache cache.ICache) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwt := r.Header.Get("Authorization")
			apiKey := r.Header.Get("X-API-Key")
			if jwt != "" {
				jwtAuth(jwtService, next, w, r)
				return
			}
			if apiKey != "" {
				apiKeyAuth(apiKeyService, apiKeyCache, next, w, r)
				return
			}
			api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		})
	}
}
