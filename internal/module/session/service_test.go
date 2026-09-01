package session_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/agent"
	appmodule "github.com/usesnipet/snipet/internal/module/app"
	session "github.com/usesnipet/snipet/internal/module/session"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	llmmocks "github.com/usesnipet/snipet/pkg/driver/llm/mocks"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/jsonx"
	"github.com/usesnipet/snipet/pkg/msg"
)

type fakeStreamIterator struct {
	events []llm.StreamEvent
	idx    int
}

func streamOf(events ...llm.StreamEvent) llm.StreamIterator {
	return &fakeStreamIterator{events: events}
}

func (it *fakeStreamIterator) Next(_ context.Context) bool {
	if it.idx >= len(it.events) {
		return false
	}
	it.idx++
	return true
}

func (it *fakeStreamIterator) Event() llm.StreamEvent { return it.events[it.idx-1] }
func (it *fakeStreamIterator) Err() error             { return nil }
func (it *fakeStreamIterator) Close() error           { return nil }

func appKeyContext(appID string) context.Context {
	return auth.SetAppKeyIdentity(context.Background(), auth.AppKeyIdentity{AppID: appID})
}

func jwtContext(userID, appCode string) context.Context {
	return auth.SetAppUserIdentity(context.Background(), auth.AppUserIdentity{UserID: userID, AppCode: appCode})
}

func newSessionService(
	t *testing.T,
	sessionRepo *mocks.MockISessionRepository,
	messageRepo *mocks.MockIExecutionMessageRepository,
	appRepo *mocks.MockIAppRepository,
	agentSvc *agent.Service,
) *session.Service {
	t.Helper()
	return session.NewService(
		sessionRepo,
		messageRepo,
		appmodule.NewService(appRepo, auth.NewAPIKeyGenerator(), auth.NewKeyHasher(), logger.NewLogger(logger.LevelError)),
		agentSvc,
	)
}

func expectAppByCode(t *testing.T, appRepo *mocks.MockIAppRepository, code, appID string) {
	t.Helper()
	appRepo.EXPECT().
		FindByCode(mock.Anything, code).
		Return(&model.App{ID: appID, Code: code}, nil)
}

func TestFindMessagesReturnsExecutionMessages(t *testing.T) {
	t.Parallel()

	appCode := "abc"
	appID := uuid.New().String()
	sessionID := uuid.New().String()
	expected := page.NewPaginated([]model.ExecutionMessage{{Message: msg.NewMessage(msg.RoleUser, "hi")}}, 1, 0, 10)

	appRepo := mocks.NewMockIAppRepository(t)
	expectAppByCode(t, appRepo, appCode, appID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		FindByIDInApp(mock.Anything, appID, sessionID, mock.Anything).
		Return(&model.Session{ID: sessionID, AppID: appID}, nil)

	messageRepo := mocks.NewMockIExecutionMessageRepository(t)
	messageRepo.EXPECT().
		FilterInSession(mock.Anything, sessionID, mock.Anything).
		Return(expected, nil)

	svc := newSessionService(t, sessionRepo, messageRepo, appRepo, nil)

	result, err := svc.FindMessages(appKeyContext(appID), appCode, sessionID, filter.Default[model.ExecutionMessage]())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindMessagesForbiddenWithoutAccess(t *testing.T) {
	t.Parallel()

	appCode := "abc"
	appID := uuid.New().String()
	sessionID := uuid.New().String()
	userID := uuid.New().String()

	appRepo := mocks.NewMockIAppRepository(t)
	expectAppByCode(t, appRepo, appCode, appID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		CheckUserAccess(mock.Anything, appID, userID, sessionID).
		Return(false, nil)

	messageRepo := mocks.NewMockIExecutionMessageRepository(t)

	svc := newSessionService(t, sessionRepo, messageRepo, appRepo, nil)

	_, err := svc.FindMessages(jwtContext(userID, appCode), appCode, sessionID, filter.Default[model.ExecutionMessage]())
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.StatusCode)
}

func TestRunDelegatesToAgentWithSessionID(t *testing.T) {
	t.Parallel()

	appCode := "abc"
	appID := uuid.New().String()
	sessionID := uuid.New().String()
	agentID := uuid.New().String()

	var created *model.Execution

	appRepo := mocks.NewMockIAppRepository(t)
	expectAppByCode(t, appRepo, appCode, appID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		FindByIDInApp(mock.Anything, appID, sessionID, mock.Anything).
		Return(&model.Session{ID: sessionID, AppID: appID, AgentID: agentID}, nil)

	agentRepo := mocks.NewMockIAgentRepository(t)
	llmID := uuid.New().String()
	agentRepo.EXPECT().
		FindByID(mock.Anything, agentID).
		Return(&model.Agent{
			ID: agentID,
			AgentToLLMs: []model.AgentToLLM{{
				AgentID:  agentID,
				LLMID:    llmID,
				Priority: 0,
				LLM: model.LLM{
					ID:            llmID,
					Provider:      "primary",
					Configuration: jsonx.JSONMap{},
				},
			}},
		}, nil)

	executionRepo := mocks.NewMockIExecutionRepository(t)
	executionRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, e *model.Execution) {
			created = e
			e.ID = uuid.New().String()
		}).
		Return(nil)
	executionRepo.EXPECT().
		UpdateByID(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	messageRepo := mocks.NewMockIExecutionMessageRepository(t)
	messageRepo.EXPECT().
		FilterInSession(mock.Anything, sessionID, mock.Anything).
		Return(page.NewPaginated([]model.ExecutionMessage{}, 0, 0, 20), nil)
	messageRepo.EXPECT().
		CreateInExecution(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	primary := llmmocks.NewMockDriver(t)
	primary.EXPECT().Info().Return(driver.Info{
		Key:                 "primary",
		Name:                "primary",
		Description:         "primary",
		ConfigurationSchema: jsonx.JSONMap{"type": "object"},
	})
	primary.EXPECT().Validate().Return(nil)
	primary.EXPECT().
		Model(mock.Anything, mock.Anything).
		Return(llm.NewModel("primary", "primary", []llm.ModelCapabilities{llm.ModelCapabilitiesToolCall}), nil)
	primary.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.Anything).
		Return(streamOf(llm.TextDeltaEvent{Text: "ok"}), nil)

	llmReg := driver.NewRegistry[llm.Driver](logger.NewLogger(logger.LevelError))
	llmReg.MustRegister(primary, nil)
	agentSvc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		runtime.NewEngine(
			manager.NewDriverManager(llmReg),
			manager.NewToolbox(manager.NewDriverManager(driver.NewRegistry[tool.Driver](logger.NewLogger(logger.LevelError)))),
			logger.NewLogger(logger.LevelError),
		),
		executionRepo,
		messageRepo,
		logger.NewLogger(logger.LevelError),
	)

	svc := newSessionService(t, sessionRepo, messageRepo, appRepo, agentSvc)

	err := svc.Run(appKeyContext(appID), appCode, sessionID, session.RunSessionDTO{Message: "hi"})
	require.NoError(t, err)

	require.NotNil(t, created)
	require.NotNil(t, created.SessionID)
	assert.Equal(t, sessionID, *created.SessionID)
	assert.Equal(t, agentID, created.AgentID)
}

func TestDeleteByIDResolvesAppCode(t *testing.T) {
	t.Parallel()

	appCode := "abc"
	appID := uuid.New().String()
	sessionID := uuid.New().String()

	appRepo := mocks.NewMockIAppRepository(t)
	expectAppByCode(t, appRepo, appCode, appID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		DeleteInApp(mock.Anything, appID, sessionID).
		Return(nil)

	svc := newSessionService(t, sessionRepo, mocks.NewMockIExecutionMessageRepository(t), appRepo, nil)

	err := svc.DeleteByID(appKeyContext(appID), appCode, sessionID)
	require.NoError(t, err)
}
