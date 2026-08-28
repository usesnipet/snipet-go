package guard

import (
	"context"
	"net/http"
	"time"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/infra/cache"
	"github.com/usesnipet/snipet/internal/module/app"
)

// cachedAppKey is what RequireAppKey caches per raw key — just enough to
// rebuild an auth.AppKeyIdentity on a cache hit without re-verifying.
type cachedAppKey struct {
	ID   string
	Code string
}

// RequireAppKey requires a valid X-App-Key header and sets an
// auth.AppKeyIdentity.
func RequireAppKey(appService *app.Service, appKeyCache cache.ICache) api.Gate {
	return func(r *http.Request) (context.Context, error) {
		key := r.Header.Get("X-App-Key")
		if key == "" {
			return nil, auth.ErrNotApplicable
		}

		if cached, found := cache.GetAs[cachedAppKey](appKeyCache, key); found {
			return auth.SetAppKeyIdentity(r.Context(), auth.AppKeyIdentity{AppID: cached.ID, Code: cached.Code}), nil
		}

		found, err := appService.VerifyKey(r.Context(), key)
		if err != nil {
			return nil, err
		}
		appKeyCache.Set(key, cachedAppKey{ID: found.ID, Code: found.Code}, cache.WithTTL(1*time.Minute))

		return auth.SetAppKeyIdentity(r.Context(), auth.AppKeyIdentity{AppID: found.ID, Code: found.Code}), nil
	}
}
