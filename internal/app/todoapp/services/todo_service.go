package services

import (
	"context"
	"errors"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/models"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoService interface {
	Create(ctx context.Context, userID primitive.ObjectID, req *models.TodoCreateRequest) (*models.Todo, error)
	GetByID(ctx context.Context, userID, todoID primitive.ObjectID) (*models.Todo, error)
	List(ctx context.Context, userID primitive.ObjectID, query *models.TodoListQuery) ([]*models.Todo, int64, error)
	Update(ctx context.Context, userID, todoID primitive.ObjectID, req *models.TodoUpdateRequest) (*models.Todo, error)
	Delete(ctx context.Context, userID, todoID primitive.ObjectID) error
}

type todoService struct {
	todoRepo repositories.TodoRepository
}

func NewTodoService(todoRepo repositories.TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

func (s *todoService) Create(ctx context.Context, userID primitive.ObjectID, req *models.TodoCreateRequest) (*models.Todo, error) {
	todo := &models.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) GetByID(ctx context.Context, userID, todoID primitive.ObjectID) (*models.Todo, error) {
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	// Ensure the todo belongs to the user
	if todo.UserID != userID {
		return nil, errors.New("unauthorized access to todo")
	}

	return todo, nil
}

func (s *todoService) List(ctx context.Context, userID primitive.ObjectID, query *models.TodoListQuery) ([]*models.Todo, int64, error) {
	// Set default values
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	todos, total, err := s.todoRepo.FindByUserID(ctx, userID, query.Page, query.PageSize, query.Completed)
	if err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (s *todoService) Update(ctx context.Context, userID, todoID primitive.ObjectID, req *models.TodoUpdateRequest) (*models.Todo, error) {
	// Check if todo exists and belongs to user
	todo, err := s.GetByID(ctx, userID, todoID)
	if err != nil {
		return nil, err
	}

	// Build update document
	update := bson.M{}
	if req.Title != nil {
		update["title"] = *req.Title
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Completed != nil {
		update["completed"] = *req.Completed
	}

	if len(update) == 0 {
		return todo, nil
	}

	if err := s.todoRepo.Update(ctx, todoID, update); err != nil {
		return nil, err
	}

	// Fetch updated todo
	updatedTodo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	return updatedTodo, nil
}

func (s *todoService) Delete(ctx context.Context, userID, todoID primitive.ObjectID) error {
	// Check if todo exists and belongs to user
	_, err := s.GetByID(ctx, userID, todoID)
	if err != nil {
		return err
	}

	return s.todoRepo.Delete(ctx, todoID)
}
