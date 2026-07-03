package knowledge_test

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
	knowledge "github.com/usesnipet/snipet/internal/module/knowledge"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

var testConfigSchema = util.JSONMap{
	"type": "object",
	"properties": util.JSONMap{
		"index": util.JSONMap{"type": "string"},
	},
	"required": []any{"index"},
}

type mockSourceDriver struct {
	mock.Mock
}

func (m *mockSourceDriver) Scan(ctx context.Context, config util.JSONMap, take *int, skip *int) ([]runtime.SourceItem, error) {
	args := m.Called(ctx, config, take, skip)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]runtime.SourceItem), args.Error(1)
}

func (m *mockSourceDriver) TestConnection(ctx context.Context, config util.JSONMap) error {
	return m.Called(ctx, config).Error(0)
}

func (m *mockSourceDriver) GetConfigurationSchema(ctx context.Context) (util.JSONMap, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(util.JSONMap), args.Error(1)
}

func newTestService(repo repository.IKnowledgeRepository, drivers map[string]*mockSourceDriver) *knowledge.Service {
	registry := runtime.NewRegistry[runtime.SourceDriver]()
	for name, driver := range drivers {
		registry.MustRegister(name, driver)
	}
	sourceManager := runtime.NewSourceManager(registry)
	return knowledge.NewService(repo, sourceManager)
}

func expectSuccessfulConnection(driver *mockSourceDriver, config util.JSONMap) {
	driver.On("GetConfigurationSchema", mock.Anything).Return(testConfigSchema, nil).Once()
	driver.On("TestConnection", mock.Anything, config).Return(nil).Once()
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

	opts := filter.Default[model.Knowledge]()
	expected := page.NewPaginated([]model.Knowledge{{Name: "Knowledge A"}}, 1, 0, 10)
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Filter(mock.Anything, opts).
		Return(expected, nil)

	svc := newTestService(repo, nil)

	result, err := svc.Filter(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	expected := &model.Knowledge{ID: id, Name: "Found"}
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		FindByID(mock.Anything, id).
		Return(expected, nil)

	svc := newTestService(repo, nil)

	result, err := svc.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresKnowledgeAndReturnsIt(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	var stored *model.Knowledge

	driver := new(mockSourceDriver)
	expectSuccessfulConnection(driver, config)

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, k *model.Knowledge) {
			stored = k
			k.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(repo, map[string]*mockSourceDriver{"pinecone": driver})

	result, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
		Name:          "Product Docs",
		Description:   "Internal documentation",
		Driver:        "pinecone",
		Configuration: config,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "Product Docs", result.Name)
	assert.Equal(t, "Internal documentation", result.Description)
	assert.Equal(t, "pinecone", result.Driver)
	assert.Equal(t, config, result.Configuration)

	require.NotNil(t, stored)
	assert.Equal(t, result.Name, stored.Name)
	assert.Equal(t, result.Description, stored.Description)
	assert.Equal(t, result.Driver, stored.Driver)
	assert.Equal(t, result.Configuration, stored.Configuration)

	driver.AssertExpectations(t)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	expectedErr := errors.New("create failed")

	driver := new(mockSourceDriver)
	expectSuccessfulConnection(driver, config)

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expectedErr)

	svc := newTestService(repo, map[string]*mockSourceDriver{"pinecone": driver})

	_, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
		Name:          "Knowledge",
		Driver:        "pinecone",
		Configuration: config,
	})
	require.ErrorIs(t, err, expectedErr)

	driver.AssertExpectations(t)
}

func TestCreateReturnsBadRequestWhenConnectionFails(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	connectionErr := errors.New("connection refused")

	driver := new(mockSourceDriver)
	driver.On("GetConfigurationSchema", mock.Anything).Return(testConfigSchema, nil).Once()
	driver.On("TestConnection", mock.Anything, config).Return(connectionErr).Once()

	repo := mocks.NewMockIKnowledgeRepository(t)

	svc := newTestService(repo, map[string]*mockSourceDriver{"pinecone": driver})

	_, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
		Name:          "Knowledge",
		Driver:        "pinecone",
		Configuration: config,
	})
	assertAppError(t, err, http.StatusBadRequest, connectionErr.Error())

	driver.AssertExpectations(t)
}

func TestUpdateDelegatesPartialFieldsToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	newName := "Updated Name"
	newDescription := "Updated description"

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, gotID string, updates *model.Knowledge) {
			assert.Equal(t, id, gotID)
			assert.Equal(t, newName, updates.Name)
			assert.Equal(t, newDescription, updates.Description)
			assert.Empty(t, updates.Driver)
			assert.Nil(t, updates.Configuration)
		}).
		Return(nil)

	svc := newTestService(repo, nil)

	err := svc.Update(context.Background(), id, knowledge.UpdateKnowledgeDTO{
		Name:        &newName,
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestUpdateDelegatesDescriptionToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	newDescription := "Updated description"

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, gotID string, updates *model.Knowledge) {
			assert.Equal(t, id, gotID)
			assert.Empty(t, updates.Name)
			assert.Equal(t, newDescription, updates.Description)
		}).
		Return(nil)

	svc := newTestService(repo, nil)

	err := svc.Update(context.Background(), id, knowledge.UpdateKnowledgeDTO{
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestDeleteByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		DeleteByID(mock.Anything, id).
		Return(nil)

	svc := newTestService(repo, nil)

	err := svc.DeleteByID(context.Background(), id)
	require.NoError(t, err)
}

func TestTestConnectionSucceedsWhenDriverValidatesAndConnects(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}

	driver := new(mockSourceDriver)
	expectSuccessfulConnection(driver, config)

	svc := newTestService(mocks.NewMockIKnowledgeRepository(t), map[string]*mockSourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", config)
	require.NoError(t, err)

	driver.AssertExpectations(t)
}

func TestTestConnectionReturnsNotFoundForUnknownDriver(t *testing.T) {
	t.Parallel()

	svc := newTestService(mocks.NewMockIKnowledgeRepository(t), nil)

	err := svc.TestConnection(context.Background(), "unknown", util.JSONMap{"index": "docs"})
	assertAppError(t, err, http.StatusNotFound, runtime.ErrSourceDriverNotFound.Error())
}

func TestTestConnectionReturnsBadRequestWhenSchemaFetchFails(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	schemaErr := errors.New("schema unavailable")

	driver := new(mockSourceDriver)
	driver.On("GetConfigurationSchema", mock.Anything).Return(nil, schemaErr).Once()

	svc := newTestService(mocks.NewMockIKnowledgeRepository(t), map[string]*mockSourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", config)
	assertAppError(t, err, http.StatusBadRequest, schemaErr.Error())

	driver.AssertExpectations(t)
}

func TestTestConnectionReturnsBadRequestWhenConfigurationInvalid(t *testing.T) {
	t.Parallel()

	driver := new(mockSourceDriver)
	driver.On("GetConfigurationSchema", mock.Anything).Return(testConfigSchema, nil).Once()

	svc := newTestService(mocks.NewMockIKnowledgeRepository(t), map[string]*mockSourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", util.JSONMap{})
	assertAppError(t, err, http.StatusBadRequest, runtime.ErrInvalidConfiguration.Error())

	driver.AssertExpectations(t)
}

func TestTestConnectionReturnsDriverConnectionError(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	connectionErr := errors.New("connection refused")

	driver := new(mockSourceDriver)
	driver.On("GetConfigurationSchema", mock.Anything).Return(testConfigSchema, nil).Once()
	driver.On("TestConnection", mock.Anything, config).Return(connectionErr).Once()

	svc := newTestService(mocks.NewMockIKnowledgeRepository(t), map[string]*mockSourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", config)
	require.ErrorIs(t, err, connectionErr)

	driver.AssertExpectations(t)
}
