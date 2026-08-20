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
	"github.com/usesnipet/snipet/internal/module/client"
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

func apiKeyContext() context.Context {
	return auth.SetApiKeyIdentity(context.Background(), auth.ApiKeyIdentity{APIKeyID: "api-key-id"})
}

func jwtContext(userID string) context.Context {
	return auth.SetClientUserIdentity(context.Background(), auth.ClientUserIdentity{UserID: userID})
}

func newSessionService(
	t *testing.T,
	sessionRepo *mocks.MockISessionRepository,
	messageRepo *mocks.MockIExecutionMessageRepository,
	clientRepo *mocks.MockIClientRepository,
	agentRepo *mocks.MockIAgentRepository,
	agentSvc *agent.Service,
) *session.Service {
	t.Helper()
	return session.NewService(
		sessionRepo,
		messageRepo,
		client.NewService(clientRepo, agentRepo, logger.NewLogger(logger.LevelError)),
		agentSvc,
	)
}

func expectClientByCode(t *testing.T, clientRepo *mocks.MockIClientRepository, code, clientID string) {
	t.Helper()
	clientRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.App{{ID: clientID, Code: code}}, 1, 0, 10), nil)
}

func TestFindMessagesReturnsExecutionMessages(t *testing.T) {
	t.Parallel()

	clientCode := "abc"
	clientID := uuid.New().String()
	sessionID := uuid.New().String()
	expected := page.NewPaginated([]model.ExecutionMessage{{Message: msg.NewMessage(msg.RoleUser, "hi")}}, 1, 0, 10)

	clientRepo := mocks.NewMockIClientRepository(t)
	expectClientByCode(t, clientRepo, clientCode, clientID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		FindByIDInClient(mock.Anything, clientID, sessionID, mock.Anything).
		Return(&model.Session{ID: sessionID, ClientID: clientID}, nil)

	messageRepo := mocks.NewMockIExecutionMessageRepository(t)
	messageRepo.EXPECT().
		FilterInSession(mock.Anything, sessionID, mock.Anything).
		Return(expected, nil)

	agentRepo := mocks.NewMockIAgentRepository(t)

	svc := newSessionService(t, sessionRepo, messageRepo, clientRepo, agentRepo, nil)

	result, err := svc.FindMessages(apiKeyContext(), clientCode, sessionID, filter.Default[model.ExecutionMessage]())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindMessagesForbiddenWithoutAccess(t *testing.T) {
	t.Parallel()

	clientCode := "abc"
	clientID := uuid.New().String()
	sessionID := uuid.New().String()
	userID := uuid.New().String()

	clientRepo := mocks.NewMockIClientRepository(t)
	expectClientByCode(t, clientRepo, clientCode, clientID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		CheckUserAccess(mock.Anything, clientID, userID, sessionID).
		Return(false, nil)

	messageRepo := mocks.NewMockIExecutionMessageRepository(t)
	agentRepo := mocks.NewMockIAgentRepository(t)

	svc := newSessionService(t, sessionRepo, messageRepo, clientRepo, agentRepo, nil)

	_, err := svc.FindMessages(jwtContext(userID), clientCode, sessionID, filter.Default[model.ExecutionMessage]())
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.StatusCode)
}

func TestRunDelegatesToAgentWithSessionID(t *testing.T) {
	t.Parallel()

	clientCode := "abc"
	clientID := uuid.New().String()
	sessionID := uuid.New().String()
	agentID := uuid.New().String()

	var created *model.Execution

	clientRepo := mocks.NewMockIClientRepository(t)
	expectClientByCode(t, clientRepo, clientCode, clientID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		FindByIDInClient(mock.Anything, clientID, sessionID, mock.Anything).
		Return(&model.Session{ID: sessionID, ClientID: clientID, AgentID: agentID}, nil)

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
		ListBySessionID(mock.Anything, sessionID).
		Return(nil, nil)
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
		Stream(mock.Anything, mock.Anything, mock.Anything).
		Return(streamOf(llm.TextDeltaEvent{Text: "ok"}), nil)

	llmReg := driver.NewRegistry[llm.Driver](logger.NewLogger(logger.LevelError))
	llmReg.MustRegister(primary, nil)
	agentSvc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		runtime.NewEngine(
			manager.NewDriver(llmReg),
			manager.NewTool(manager.NewDriver(driver.NewRegistry[tool.Driver](logger.NewLogger(logger.LevelError)))),
			logger.NewLogger(logger.LevelError),
		),
		executionRepo,
		messageRepo,
		logger.NewLogger(logger.LevelError),
	)

	svc := newSessionService(t, sessionRepo, messageRepo, clientRepo, agentRepo, agentSvc)

	err := svc.Run(apiKeyContext(), clientCode, sessionID, session.RunSessionDTO{Message: "hi"})
	require.NoError(t, err)

	require.NotNil(t, created)
	require.NotNil(t, created.SessionID)
	assert.Equal(t, sessionID, *created.SessionID)
	assert.Equal(t, agentID, created.AgentID)
}

func TestDeleteByIDResolvesClientCode(t *testing.T) {
	t.Parallel()

	clientCode := "abc"
	clientID := uuid.New().String()
	sessionID := uuid.New().String()

	clientRepo := mocks.NewMockIClientRepository(t)
	expectClientByCode(t, clientRepo, clientCode, clientID)

	sessionRepo := mocks.NewMockISessionRepository(t)
	sessionRepo.EXPECT().
		DeleteInClient(mock.Anything, clientID, sessionID).
		Return(nil)

	agentRepo := mocks.NewMockIAgentRepository(t)

	svc := newSessionService(t, sessionRepo, mocks.NewMockIExecutionMessageRepository(t), clientRepo, agentRepo, nil)

	err := svc.DeleteByID(apiKeyContext(), clientCode, sessionID)
	require.NoError(t, err)
}
