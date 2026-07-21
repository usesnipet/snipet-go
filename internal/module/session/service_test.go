package session_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
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
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
)

type sessionTestLLM struct {
	key string
}

func (s *sessionTestLLM) Info() driver.Info {
	return driver.Info{
		Name:                s.key,
		Description:         s.key,
		ConfigurationSchema: util.JSONMap{"type": "object"},
	}
}
func (s *sessionTestLLM) TestConnection(context.Context, util.JSONMap) error { return nil }
func (s *sessionTestLLM) Generate(
	context.Context,
	util.JSONMap,
	string,
	[]message.Message,
) (message.Message, error) {
	return message.Message{Role: message.MessageRoleFinal, Content: "ok"}, nil
}

func apiKeyContext() context.Context {
	id := "api-key-id"
	return auth.SetPrincipal(context.Background(), auth.NewPrincipal(auth.PrincipalTypeAPIKey, &id, nil))
}

func jwtContext(userID string) context.Context {
	return auth.SetPrincipal(context.Background(), auth.NewPrincipal(
		auth.PrincipalTypeJWT,
		nil,
		&auth.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: userID},
		},
	))
}

func newSessionService(
	t *testing.T,
	sessionRepo *mocks.MockISessionRepository,
	messageRepo *mocks.MockIExecutionMessageRepository,
	clientRepo *mocks.MockIClientRepository,
	agentSvc *agent.Service,
) *session.Service {
	t.Helper()
	return session.NewService(
		sessionRepo,
		messageRepo,
		client.NewService(clientRepo),
		agentSvc,
	)
}

func expectClientByCode(t *testing.T, clientRepo *mocks.MockIClientRepository, code, clientID string) {
	t.Helper()
	clientRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.Client{{ID: clientID, Code: code}}, 1, 0, 10), nil)
}

func TestFindMessagesReturnsExecutionMessages(t *testing.T) {
	t.Parallel()

	clientCode := "abc"
	clientID := uuid.New().String()
	sessionID := uuid.New().String()
	expected := page.NewPaginated([]model.ExecutionMessage{{Content: "hi"}}, 1, 0, 10)

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

	svc := newSessionService(t, sessionRepo, messageRepo, clientRepo, nil)

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
	svc := newSessionService(t, sessionRepo, messageRepo, clientRepo, nil)

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
	agentRepo.EXPECT().
		FindByID(mock.Anything, agentID).
		Return(&model.Agent{
			ID: agentID,
			Configuration: model.AgentConfiguration{
				LLMs:  []runtime.LLMConfig{{Key: "primary", Config: util.JSONMap{}}},
				Tools: runtime.ToolConfig{},
			},
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

	llmReg := registry.New[driver.ILLM]()
	llmReg.MustRegister("primary", &sessionTestLLM{key: "primary"})
	agentSvc := agent.NewService(
		agentRepo,
		runtime.NewEngine(
			driver.NewManager(registry.New[driver.ITool]()),
			driver.NewManager(llmReg),
			logger.NewLogger(logger.LevelError),
		),
		driver.NewManager(llmReg),
		driver.NewManager(registry.New[driver.ITool]()),
		executionRepo,
		messageRepo,
		logger.NewLogger(logger.LevelError),
	)

	svc := newSessionService(t, sessionRepo, messageRepo, clientRepo, agentSvc)

	err := svc.Run(apiKeyContext(), clientCode, sessionID, session.RunSessionDTO{Message: "hi"}, nil)
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

	svc := newSessionService(t, sessionRepo, mocks.NewMockIExecutionMessageRepository(t), clientRepo, nil)

	err := svc.DeleteByID(apiKeyContext(), clientCode, sessionID)
	require.NoError(t, err)
}
