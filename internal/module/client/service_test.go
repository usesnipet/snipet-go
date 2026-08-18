package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	client "github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/page"
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
	assert.Equal(t, statusCode, appErr.StatusCode)
}

func TestFindByCodeInTenantRejectsCrossTenant(t *testing.T) {
	t.Parallel()

	clientRepo := mocks.NewMockIClientRepository(t)
	clientRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.Client{Code: "acme", TenantID: "other-tenant"}, nil)

	svc := client.NewService(clientRepo, mocks.NewMockIAgentRepository(t), logger.NewLogger(logger.LevelError))

	_, err := svc.FindByCodeInTenant(adminCtx(), tenantID, "acme")
	assertAppError(t, err, http.StatusNotFound)
}

func TestGetAgentsScopesToClientTenant(t *testing.T) {
	t.Parallel()

	clientRepo := mocks.NewMockIClientRepository(t)
	clientRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.Client{{Code: "acme", TenantID: tenantID}}, 1, 0, 1), nil)

	expected := page.NewPaginated([]model.Agent{{Name: "Agent A", TenantID: tenantID}}, 1, 0, 10)
	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Filter(mock.Anything, mock.MatchedBy(func(opts *filter.Options[model.Agent]) bool {
			where, ok := opts.Where.Fields["tenant_id"]
			return ok && len(where.Value) == 1 && where.Value[0] == tenantID
		})).
		Return(expected, nil)

	svc := client.NewService(clientRepo, agentRepo, logger.NewLogger(logger.LevelError))

	result, err := svc.GetAgents(context.Background(), "acme", filter.Default[model.Agent]())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
