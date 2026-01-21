package grpc

import (
	"context"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/todo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
	todoUsecase usecase.TodoUsecase
}

func NewTodoServer(todoUsecase usecase.TodoUsecase) *TodoServer {
	return &TodoServer{
		todoUsecase: todoUsecase,
	}
}

func (s *TodoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.TodoResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	todoReq := &dto.TodoCreateRequest{
		Title:       req.Title,
		Description: req.Description,
	}

	todo, err := s.todoUsecase.Create(ctx, userID, todoReq)
	if err != nil {
		return nil, err
	}

	return &pb.TodoResponse{
		Id:          todo.ID,
		UserId:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}

func (s *TodoServer) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.TodoResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	todoID, err := primitive.ObjectIDFromHex(req.TodoId)
	if err != nil {
		return nil, err
	}

	todo, err := s.todoUsecase.GetByID(ctx, userID, todoID)
	if err != nil {
		return nil, err
	}

	return &pb.TodoResponse{
		Id:          todo.ID,
		UserId:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}

func (s *TodoServer) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	var completed *bool
	if req.Completed != nil {
		c := *req.Completed
		completed = &c
	}

	query := &dto.TodoListQuery{
		Page:      int(req.Page),
		PageSize:  int(req.PageSize),
		Completed: completed,
	}

	todos, total, err := s.todoUsecase.List(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	pbTodos := make([]*pb.TodoResponse, len(todos))
	for i, todo := range todos {
		pbTodos[i] = &pb.TodoResponse{
			Id:          todo.ID,
			UserId:      todo.UserID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	return &pb.ListTodosResponse{
		Todos:    pbTodos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *TodoServer) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.TodoResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	todoID, err := primitive.ObjectIDFromHex(req.TodoId)
	if err != nil {
		return nil, err
	}

	updateReq := &dto.TodoUpdateRequest{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	}

	todo, err := s.todoUsecase.Update(ctx, userID, todoID, updateReq)
	if err != nil {
		return nil, err
	}

	return &pb.TodoResponse{
		Id:          todo.ID,
		UserId:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}

func (s *TodoServer) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	todoID, err := primitive.ObjectIDFromHex(req.TodoId)
	if err != nil {
		return nil, err
	}

	if err := s.todoUsecase.Delete(ctx, userID, todoID); err != nil {
		return &pb.DeleteTodoResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteTodoResponse{
		Success: true,
		Message: "Todo deleted successfully",
	}, nil
}
