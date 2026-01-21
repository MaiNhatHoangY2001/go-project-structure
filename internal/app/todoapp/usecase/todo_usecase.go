package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoUsecase interface {
	Create(ctx context.Context, userID primitive.ObjectID, req *dto.TodoCreateRequest) (*dto.TodoResponse, error)
	GetByID(ctx context.Context, userID, todoID primitive.ObjectID) (*dto.TodoResponse, error)
	List(ctx context.Context, userID primitive.ObjectID, query *dto.TodoListQuery) ([]*dto.TodoResponse, int64, error)
	Update(ctx context.Context, userID, todoID primitive.ObjectID, req *dto.TodoUpdateRequest) (*dto.TodoResponse, error)
	Delete(ctx context.Context, userID, todoID primitive.ObjectID) error
}

type todoUsecase struct {
	todoRepo domain.TodoRepository
}

func NewTodoUsecase(todoRepo domain.TodoRepository) TodoUsecase {
	return &todoUsecase{
		todoRepo: todoRepo,
	}
}

func (u *todoUsecase) Create(ctx context.Context, userID primitive.ObjectID, req *dto.TodoCreateRequest) (*dto.TodoResponse, error) {
	todo := &entity.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	if err := u.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return &dto.TodoResponse{
		ID:          todo.ID.Hex(),
		UserID:      todo.UserID.Hex(),
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *todoUsecase) GetByID(ctx context.Context, userID, todoID primitive.ObjectID) (*dto.TodoResponse, error) {
	todo, err := u.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	if todo.UserID != userID {
		return nil, errors.New("unauthorized access to todo")
	}

	return &dto.TodoResponse{
		ID:          todo.ID.Hex(),
		UserID:      todo.UserID.Hex(),
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *todoUsecase) List(ctx context.Context, userID primitive.ObjectID, query *dto.TodoListQuery) ([]*dto.TodoResponse, int64, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	todos, total, err := u.todoRepo.FindByUserID(ctx, userID, query.Page, query.PageSize, query.Completed)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = &dto.TodoResponse{
			ID:          todo.ID.Hex(),
			UserID:      todo.UserID.Hex(),
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
		}
	}

	return responses, total, nil
}

func (u *todoUsecase) Update(ctx context.Context, userID, todoID primitive.ObjectID, req *dto.TodoUpdateRequest) (*dto.TodoResponse, error) {
	todo, err := u.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	if todo.UserID != userID {
		return nil, errors.New("unauthorized access to todo")
	}

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
		return &dto.TodoResponse{
			ID:          todo.ID.Hex(),
			UserID:      todo.UserID.Hex(),
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
		}, nil
	}

	if err := u.todoRepo.Update(ctx, todoID, update); err != nil {
		return nil, err
	}

	updatedTodo, err := u.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	return &dto.TodoResponse{
		ID:          updatedTodo.ID.Hex(),
		UserID:      updatedTodo.UserID.Hex(),
		Title:       updatedTodo.Title,
		Description: updatedTodo.Description,
		Completed:   updatedTodo.Completed,
		CreatedAt:   updatedTodo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedTodo.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *todoUsecase) Delete(ctx context.Context, userID, todoID primitive.ObjectID) error {
	todo, err := u.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("todo not found")
		}
		return err
	}

	if todo.UserID != userID {
		return errors.New("unauthorized access to todo")
	}

	return u.todoRepo.Delete(ctx, todoID)
}
