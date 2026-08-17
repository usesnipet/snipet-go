package llm_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	llmmodule "github.com/usesnipet/snipet/internal/module/llm"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver"
	llmdriver "github.com/usesnipet/snipet/pkg/driver/llm"
	llmdrivermocks "github.com/usesnipet/snipet/pkg/driver/llm/mocks"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

const (
	tenantID      = "11111111-1111-1111-1111-111111111111"
	otherTenantID = "99999999-9999-9999-9999-999999999999"
	userID        = "22222222-2222-2222-2222-222222222222"
)

var permissiveSchema = jsonx.JSONMap{"type": "object"}

func regularUser() *model.User {
	return &model.User{ID: userID, Name: "Regular", Email: "user@example.com"}
}

// ctxFor mirrors member/service_test.go's helper — builds the context
// guard.RequireUser would have populated for user, with the given tenant
// memberships.
func ctxFor(user *model.User, memberships ...*model.Member) context.Context {
	return auth.SetUserIdentity(context.Background(), &auth.UserIdentity{User: user, Memberships: memberships})
}

func assertAppError(t *testing.T, err error, statusCode int) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, statusCode, appErr.StatusCode)
}

func newTestService(t *testing.T, repo repository.ILLMRepository) *llmmodule.Service {
	t.Helper()

	reg := driver.NewRegistry[llmdriver.Driver](logger.NewLogger(logger.LevelError))
	fake := llmdrivermocks.NewMockDriver(t)
	fake.EXPECT().Info().Return(driver.Info{Key: "openai", ConfigurationSchema: permissiveSchema}).Maybe()
	fake.EXPECT().Validate().Return(nil).Maybe()
	reg.MustRegister(fake, nil)
	llmManager := manager.NewDriver(reg)

	return llmmodule.NewService(repo, llmManager)
}

func TestFilterRejectsNonMember(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, mocks.NewMockILLMRepository(t))

	_, err := svc.Filter(ctxFor(regularUser()), tenantID, nil)
	assertAppError(t, err, 403)
}

func TestFindByIDRejectsCrossTenant(t *testing.T) {
	t.Parallel()

	llmRepo := mocks.NewMockILLMRepository(t)
	llmRepo.EXPECT().
		FindByID(mock.Anything, "llm-1").
		Return(&model.LLM{ID: "llm-1", TenantID: otherTenantID}, nil)

	svc := newTestService(t, llmRepo)

	ctx := ctxFor(regularUser(), &model.Member{ID: "m1", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true})
	_, err := svc.FindByID(ctx, tenantID, "llm-1")
	assertAppError(t, err, 404)
}

func TestCreateStampsTenantID(t *testing.T) {
	t.Parallel()

	llmRepo := mocks.NewMockILLMRepository(t)
	llmRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(l *model.LLM) bool { return l.TenantID == tenantID })).
		Return(nil)

	svc := newTestService(t, llmRepo)

	ctx := ctxFor(regularUser(), &model.Member{ID: "m1", UserID: userID, TenantID: tenantID, Role: model.RoleUser, IsActive: true})
	created, err := svc.Create(ctx, tenantID, llmmodule.CreateLLMDTO{
		Name:          "GPT",
		Provider:      "openai",
		Configuration: jsonx.JSONMap{},
	})
	require.NoError(t, err)
	assert.Equal(t, tenantID, created.TenantID)
}
