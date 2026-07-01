package bot_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	bot "github.com/usesnipet/snipet/internal/module/bot"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
)

func newTestService(
	botRepo repository.IBotRepository,
	clientRepo repository.IClientRepository,
) *bot.Service {
	return bot.NewService(botRepo, clientRepo)
}

func assertAppError(t *testing.T, err error, statusCode int, message string) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, statusCode, appErr.StatusCode)
	assert.Equal(t, message, appErr.Message)
}

func TestFilterDelegatesToRepository(t *testing.T) {
	t.Parallel()

	expected := page.NewPaginated([]model.Bot{{Name: "Bot A"}}, 1, 0, 10)
	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Run(func(_ context.Context, opts *filter.Options[model.Bot]) {
			assert.Equal(t, filter.Default[model.Bot]().Take, opts.Take)
		}).
		Return(expected, nil)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	result, err := svc.Filter(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	expected := &model.Bot{ID: uuid.MustParse(id), Name: "Found"}
	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		FindByID(mock.Anything, id).
		Return(expected, nil)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	result, err := svc.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresBotAndReturnsIt(t *testing.T) {
	t.Parallel()

	config := model.BotConfiguration{LLMs: []any{"gpt-4"}}
	var stored *model.Bot

	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, b *model.Bot) {
			stored = b
			b.ID = uuid.New()
		}).
		Return(nil)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	result, err := svc.Create(context.Background(), bot.CreateBotDTO{
		Name:          "Support Bot",
		Description:   "Handles support tickets",
		Configuration: config,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "Support Bot", result.Name)
	assert.Equal(t, "Handles support tickets", result.Description)
	assert.Equal(t, config, result.Configuration)

	require.NotNil(t, stored)
	assert.Equal(t, result.Name, stored.Name)
	assert.Equal(t, result.Description, stored.Description)
	assert.Equal(t, result.Configuration, stored.Configuration)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("create failed")
	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expectedErr)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	_, err := svc.Create(context.Background(), bot.CreateBotDTO{
		Name:          "Bot",
		Configuration: model.BotConfiguration{},
	})
	require.ErrorIs(t, err, expectedErr)
}

func TestUpdateDelegatesPartialFieldsToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	newName := "Updated Name"
	newDescription := "Updated description"

	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, gotID string, updates *model.Bot) {
			assert.Equal(t, id, gotID)
			assert.Equal(t, newName, updates.Name)
			assert.Equal(t, newDescription, updates.Description)
			assert.Empty(t, updates.Configuration.LLMs)
		}).
		Return(nil)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	err := svc.Update(context.Background(), id, bot.UpdateBotDTO{
		Name:        &newName,
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestUpdateDelegatesConfigurationToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	config := model.BotConfiguration{LLMs: []any{"claude"}}

	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, _ string, updates *model.Bot) {
			assert.Equal(t, config, updates.Configuration)
		}).
		Return(nil)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	err := svc.Update(context.Background(), id, bot.UpdateBotDTO{
		Configuration: &config,
	})
	require.NoError(t, err)
}

func TestDeleteByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		DeleteByID(mock.Anything, id).
		Return(nil)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	err := svc.DeleteByID(context.Background(), id)
	require.NoError(t, err)
}

func TestLinkClientToBotRejectsInvalidBotID(t *testing.T) {
	t.Parallel()

	svc := newTestService(
		mocks.NewMockIBotRepository(t),
		mocks.NewMockIClientRepository(t),
	)

	err := svc.LinkClientToBot(context.Background(), bot.LinkClientToBotDTO{
		ClientCode: "CLIENT01",
		BotID:      "not-a-uuid",
	})
	assertAppError(t, err, http.StatusBadRequest, "invalid bot id")
}

func TestLinkClientToBotReturnsClientLookupError(t *testing.T) {
	t.Parallel()

	botID := uuid.New().String()
	expectedErr := apperr.NotFound("client not found")

	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		FindByID(mock.Anything, botID).
		Return(&model.Bot{ID: uuid.MustParse(botID)}, nil)

	clientRepo := mocks.NewMockIClientRepository(t)
	clientRepo.EXPECT().
		FindByCode(mock.Anything, "CLIENT01").
		Return(nil, expectedErr)

	svc := newTestService(botRepo, clientRepo)

	err := svc.LinkClientToBot(context.Background(), bot.LinkClientToBotDTO{
		ClientCode: "CLIENT01",
		BotID:      botID,
	})
	assertAppError(t, err, http.StatusNotFound, "client not found")
}

func TestLinkClientToBotReturnsBotLookupError(t *testing.T) {
	t.Parallel()

	expectedErr := apperr.NotFound("entity not found")
	botID := uuid.New().String()

	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		FindByID(mock.Anything, botID).
		Return(nil, expectedErr)

	svc := newTestService(botRepo, mocks.NewMockIClientRepository(t))

	err := svc.LinkClientToBot(context.Background(), bot.LinkClientToBotDTO{
		ClientCode: "CLIENT01",
		BotID:      botID,
	})
	assertAppError(t, err, http.StatusNotFound, "entity not found")
}

func TestLinkClientToBotLinksClientToBot(t *testing.T) {
	t.Parallel()

	clientID := uuid.New()
	botID := uuid.New()
	client := &model.Client{ID: clientID, Code: "CLIENT01"}

	clientRepo := mocks.NewMockIClientRepository(t)
	clientRepo.EXPECT().
		FindByCode(mock.Anything, "CLIENT01").
		Return(client, nil)

	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		FindByID(mock.Anything, botID.String()).
		Return(&model.Bot{ID: botID}, nil)
	botRepo.EXPECT().
		LinkBotToClient(mock.Anything, clientID, botID).
		Return(nil)

	svc := newTestService(botRepo, clientRepo)

	err := svc.LinkClientToBot(context.Background(), bot.LinkClientToBotDTO{
		ClientCode: "CLIENT01",
		BotID:      botID.String(),
	})
	require.NoError(t, err)
}

func TestLinkClientToBotReturnsLinkError(t *testing.T) {
	t.Parallel()

	clientID := uuid.New()
	botID := uuid.New()
	expectedErr := errors.New("link failed")

	clientRepo := mocks.NewMockIClientRepository(t)
	clientRepo.EXPECT().
		FindByCode(mock.Anything, "CLIENT01").
		Return(&model.Client{ID: clientID, Code: "CLIENT01"}, nil)

	botRepo := mocks.NewMockIBotRepository(t)
	botRepo.EXPECT().
		FindByID(mock.Anything, botID.String()).
		Return(&model.Bot{ID: botID}, nil)
	botRepo.EXPECT().
		LinkBotToClient(mock.Anything, clientID, botID).
		Return(expectedErr)

	svc := newTestService(botRepo, clientRepo)

	err := svc.LinkClientToBot(context.Background(), bot.LinkClientToBotDTO{
		ClientCode: "CLIENT01",
		BotID:      botID.String(),
	})
	require.ErrorIs(t, err, expectedErr)
}
