package clients

import (
	"context"
	"fmt"

	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TodoClient struct {
	client pb.TodoServiceClient
	conn   *grpc.ClientConn
}

func NewTodoClient(host string, port string) (*TodoClient, error) {
	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to todo service: %w", err)
	}

	return &TodoClient{
		client: pb.NewTodoServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *TodoClient) CreateTodo(ctx context.Context, userId, title, description string) (*pb.TodoResponse, error) {
	req := &pb.CreateTodoRequest{
		UserId:      userId,
		Title:       title,
		Description: description,
	}

	return c.client.CreateTodo(ctx, req)
}

func (c *TodoClient) GetTodo(ctx context.Context, userId, todoId string) (*pb.TodoResponse, error) {
	req := &pb.GetTodoRequest{
		UserId: userId,
		TodoId: todoId,
	}

	return c.client.GetTodo(ctx, req)
}

func (c *TodoClient) ListTodos(ctx context.Context, userId string, page, pageSize int32) ([]*pb.TodoResponse, int64, error) {
	req := &pb.ListTodosRequest{
		UserId:   userId,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListTodos(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Todos, resp.Total, nil
}

func (c *TodoClient) UpdateTodo(ctx context.Context, userId, todoId, title, description string, completed bool) (*pb.TodoResponse, error) {
	req := &pb.UpdateTodoRequest{
		UserId:      userId,
		TodoId:      todoId,
		Title:       &title,
		Description: &description,
		Completed:   &completed,
	}

	return c.client.UpdateTodo(ctx, req)
}

func (c *TodoClient) DeleteTodo(ctx context.Context, userId, todoId string) error {
	req := &pb.DeleteTodoRequest{
		UserId: userId,
		TodoId: todoId,
	}

	_, err := c.client.DeleteTodo(ctx, req)
	return err
}

func (c *TodoClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
