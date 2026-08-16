package guard

import (
	"context"
	"net/http"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
)

// RequireUser requires a valid platform-staff JWT (Authorization: Bearer)
// and resolves it into an auth.UserIdentity — the user's profile plus every
// tenant membership they hold — in the same pass, so handlers/services never
// need to re-query for it.
func RequireUser(
	jwtService *auth.JWTService[*auth.PlatformUserClaims],
	userRepo repository.IUserRepository,
	memberRepo repository.IMemberRepository,
) Gate {
	return func(r *http.Request) (context.Context, error) {
		claims, err := verifyBearerJWT(r, jwtService)
		if err != nil {
			return nil, err
		}
		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			return nil, apperr.Unauthorized("unauthorized")
		}

		ctx := r.Context()

		user, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			return nil, apperr.Unauthorized("unauthorized")
		}

		memberships, err := memberRepo.Filter(ctx, filter.New[model.Member](
			filter.WhereEq("user_id", userID),
			filter.Take(1000),
		))
		if err != nil {
			return nil, err
		}
		data := make([]*model.Member, 0, len(memberships.Data))
		for i := range memberships.Data {
			data = append(data, &memberships.Data[i])
		}

		return auth.SetUserIdentity(ctx, &auth.UserIdentity{User: user, Memberships: data}), nil
	}
}
