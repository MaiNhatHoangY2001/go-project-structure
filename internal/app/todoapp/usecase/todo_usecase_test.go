package usecase

import (
	"context"
	"testing"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}

func (m *MockTodoRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int, completed *bool) ([]*entity.Todo, int64, error) {
	args := m.Called(ctx, userID, page, pageSize, completed)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Todo), args.Get(1).(int64), args.Error(2)
}

func (m *MockTodoRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTodoCreate_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	req := &dto.TodoCreateRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Todo")).Return(nil)

	todo, err := todoUsecase.Create(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, req.Title, todo.Title)
	assert.Equal(t, req.Description, todo.Description)
	assert.Equal(t, userID.Hex(), todo.UserID)
	assert.False(t, todo.Completed)
	mockRepo.AssertExpectations(t)
}

func TestTodoGetByID_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	todoID := primitive.NewObjectID()

	expectedTodo := &entity.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
	}

	mockRepo.On("FindByID", mock.Anything, todoID).Return(expectedTodo, nil)

	todo, err := todoUsecase.GetByID(context.Background(), userID, todoID)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, expectedTodo.ID.Hex(), todo.ID)
	assert.Equal(t, expectedTodo.Title, todo.Title)
	mockRepo.AssertExpectations(t)
}

func TestTodoGetByID_NotFound(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	todoID := primitive.NewObjectID()

	mockRepo.On("FindByID", mock.Anything, todoID).Return(nil, mongo.ErrNoDocuments)

	todo, err := todoUsecase.GetByID(context.Background(), userID, todoID)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Equal(t, "todo not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTodoGetByID_UnauthorizedAccess(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	anotherUserID := primitive.NewObjectID()
	todoID := primitive.NewObjectID()

	expectedTodo := &entity.Todo{
		ID:          todoID,
		UserID:      anotherUserID,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
	}

	mockRepo.On("FindByID", mock.Anything, todoID).Return(expectedTodo, nil)

	todo, err := todoUsecase.GetByID(context.Background(), userID, todoID)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Equal(t, "unauthorized access to todo", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTodoList_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	query := &dto.TodoListQuery{
		Page:     1,
		PageSize: 10,
	}

	expectedTodos := []*entity.Todo{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Title:       "Todo 1",
			Description: "Description 1",
			Completed:   false,
		},
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Title:       "Todo 2",
			Description: "Description 2",
			Completed:   true,
		},
	}

	mockRepo.On("FindByUserID", mock.Anything, userID, 1, 10, (*bool)(nil)).Return(expectedTodos, int64(2), nil)

	todos, total, err := todoUsecase.List(context.Background(), userID, query)

	assert.NoError(t, err)
	assert.NotNil(t, todos)
	assert.Equal(t, 2, len(todos))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTodoUpdate_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	todoID := primitive.NewObjectID()

	existingTodo := &entity.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Old Title",
		Description: "Old Description",
		Completed:   false,
	}

	newTitle := "New Title"
	completed := true
	req := &dto.TodoUpdateRequest{
		Title:     &newTitle,
		Completed: &completed,
	}

	updatedTodo := &entity.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       newTitle,
		Description: "Old Description",
		Completed:   completed,
	}

	mockRepo.On("FindByID", mock.Anything, todoID).Return(existingTodo, nil).Once()
	mockRepo.On("Update", mock.Anything, todoID, mock.AnythingOfType("primitive.M")).Return(nil)
	mockRepo.On("FindByID", mock.Anything, todoID).Return(updatedTodo, nil).Once()

	todo, err := todoUsecase.Update(context.Background(), userID, todoID, req)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, newTitle, todo.Title)
	assert.Equal(t, completed, todo.Completed)
	mockRepo.AssertExpectations(t)
}

func TestTodoDelete_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	todoID := primitive.NewObjectID()

	existingTodo := &entity.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
	}

	mockRepo.On("FindByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockRepo.On("Delete", mock.Anything, todoID).Return(nil)

	err := todoUsecase.Delete(context.Background(), userID, todoID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTodoDelete_NotFound(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoUsecase := NewTodoUsecase(mockRepo)

	userID := primitive.NewObjectID()
	todoID := primitive.NewObjectID()

	mockRepo.On("FindByID", mock.Anything, todoID).Return(nil, mongo.ErrNoDocuments)

	err := todoUsecase.Delete(context.Background(), userID, todoID)

	assert.Error(t, err)
	assert.Equal(t, "todo not found", err.Error())
	mockRepo.AssertExpectations(t)
}
