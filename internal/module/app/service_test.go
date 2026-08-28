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

func assertAppError(t *testing.T, err error, statusCode int) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, statusCode, appErr.StatusCode)
}

func TestFindByCodeDelegatesToRepository(t *testing.T) {
	t.Parallel()

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{Code: "acme"}, nil)

	svc := appmodule.NewService(appRepo, auth.NewAPIKeyGenerator(), auth.NewKeyHasher(), logger.NewLogger(logger.LevelError))

	data, err := svc.FindByCode(context.Background(), "acme")
	require.NoError(t, err)
	require.Equal(t, "acme", data.Code)
}

func TestFindPublicByCodeReturnsPublicFields(t *testing.T) {
	t.Parallel()

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{Code: "acme", Name: "Acme", Description: "desc", Public: true}, nil)

	svc := appmodule.NewService(appRepo, auth.NewAPIKeyGenerator(), auth.NewKeyHasher(), logger.NewLogger(logger.LevelError))

	data, err := svc.FindPublicByCode(context.Background(), "acme")
	require.NoError(t, err)
	require.Equal(t, &appmodule.PublicAppDTO{Code: "acme", Name: "Acme", Description: "desc"}, data)
}

func TestFindPublicByCodeRejectsNonPublicApp(t *testing.T) {
	t.Parallel()

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{Code: "acme", Public: false}, nil)

	svc := appmodule.NewService(appRepo, auth.NewAPIKeyGenerator(), auth.NewKeyHasher(), logger.NewLogger(logger.LevelError))

	_, err := svc.FindPublicByCode(context.Background(), "acme")
	assertAppError(t, err, http.StatusNotFound)
}
