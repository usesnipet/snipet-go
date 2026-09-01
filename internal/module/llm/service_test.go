package llm_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	llmmodule "github.com/usesnipet/snipet/internal/module/llm"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver"
	llmdriver "github.com/usesnipet/snipet/pkg/driver/llm"
	llmdrivermocks "github.com/usesnipet/snipet/pkg/driver/llm/mocks"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

var permissiveSchema = jsonx.JSONMap{"type": "object"}

func newTestService(t *testing.T, repo *mocks.MockILLMRepository) *llmmodule.Service {
	t.Helper()

	reg := driver.NewRegistry[llmdriver.Driver](logger.NewLogger(logger.LevelError))
	fake := llmdrivermocks.NewMockDriver(t)
	fake.EXPECT().Info().Return(driver.Info{Key: "openai", ConfigurationSchema: permissiveSchema}).Maybe()
	fake.EXPECT().Validate().Return(nil).Maybe()
	reg.MustRegister(fake, nil)
	llmManager := manager.NewDriverManager(reg)

	return llmmodule.NewService(repo, llmManager)
}

func TestFindByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	llmRepo := mocks.NewMockILLMRepository(t)
	llmRepo.EXPECT().
		FindByID(mock.Anything, "llm-1").
		Return(&model.LLM{ID: "llm-1", Name: "Found"}, nil)

	svc := newTestService(t, llmRepo)

	result, err := svc.FindByID(context.Background(), "llm-1")
	require.NoError(t, err)
	assert.Equal(t, "Found", result.Name)
}

func TestCreateStoresLLM(t *testing.T) {
	t.Parallel()

	llmRepo := mocks.NewMockILLMRepository(t)
	llmRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(l *model.LLM) bool { return l.Name == "GPT" })).
		Return(nil)

	svc := newTestService(t, llmRepo)

	created, err := svc.Create(context.Background(), llmmodule.CreateLLMDTO{
		Name:          "GPT",
		Provider:      "openai",
		Configuration: jsonx.JSONMap{},
	})
	require.NoError(t, err)
	assert.Equal(t, "GPT", created.Name)
	assert.Equal(t, "openai", created.Provider)
}
