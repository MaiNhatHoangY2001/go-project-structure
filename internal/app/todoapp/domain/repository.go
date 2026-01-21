package domain

import (
	"context"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error)
}

// TodoRepository defines the interface for todo data access
type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Todo, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int, completed *bool) ([]*entity.Todo, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
