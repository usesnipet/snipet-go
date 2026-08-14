package middleware

import (
	"errors"
	"net/http"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
)

// IdentityMiddleware loads the authenticated platform user's profile and
// every tenant membership they hold, and stashes it on the request context
// as an auth.Identity (see auth.CurrentIdentity) — so handlers and services
// read it directly instead of each re-querying the user/member
// repositories. Must run after RequirePlatformJWT, which is what puts the
// platform user id on the context in the first place.
type IdentityMiddleware struct {
	userRepo   repository.IUserRepository
	memberRepo repository.IMemberRepository
}

func RequireIdentity(userRepo repository.IUserRepository, memberRepo repository.IMemberRepository) IdentityMiddleware {
	return IdentityMiddleware{userRepo: userRepo, memberRepo: memberRepo}
}

func (m IdentityMiddleware) Handler() api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userID, err := auth.PlatformUserID(ctx)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}

			user, err := m.userRepo.FindByID(ctx, userID)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}

			memberships, err := m.memberRepo.Filter(ctx, filter.New[model.Member](
				filter.WhereEq("user_id", userID),
				filter.Take(1000),
			))
			if err != nil {
				api.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			data := make([]*model.Member, 0, len(memberships.Data))
			for i := range memberships.Data {
				data = append(data, &memberships.Data[i])
			}

			ctx = auth.SetIdentity(ctx, &auth.Identity{User: user, Memberships: data})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
