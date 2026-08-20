package app_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	appmodule "github.com/usesnipet/snipet/internal/module/app"
	"github.com/usesnipet/snipet/internal/repository/mocks"
)

const (
	tenantID = "11111111-1111-1111-1111-111111111111"
	adminID  = "22222222-2222-2222-2222-222222222222"
)

func adminCtx() context.Context {
	admin := &model.User{ID: adminID, Name: "Admin", Email: "admin@example.com"}
	return auth.SetUserIdentity(context.Background(), &auth.UserIdentity{
		User:        admin,
		Memberships: []*model.Member{{ID: "m1", UserID: adminID, TenantID: tenantID, Role: model.RoleAdmin, IsActive: true}},
	})
}

func assertAppError(t *testing.T, err error, statusCode int) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, statusCode, appErr.StatusCode)
}

func TestFindByCodeInTenantRejectsCrossTenant(t *testing.T) {
	t.Parallel()

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{Code: "acme", TenantID: "other-tenant"}, nil)

	svc := appmodule.NewService(appRepo, auth.NewAPIKeyGenerator(), auth.NewKeyHasher(), logger.NewLogger(logger.LevelError))

	_, err := svc.FindByCodeInTenant(adminCtx(), tenantID, "acme")
	assertAppError(t, err, http.StatusNotFound)
}
