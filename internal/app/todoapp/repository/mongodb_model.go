package repository

import (
"time"

"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserModel represents the MongoDB document structure for users
type UserModel struct {
ID        primitive.ObjectID `bson:"_id,omitempty"`
Email     string             `bson:"email"`
Password  string             `bson:"password"`
Name      string             `bson:"name"`
CreatedAt time.Time          `bson:"created_at"`
UpdatedAt time.Time          `bson:"updated_at"`
}

// ToEntity converts UserModel to domain entity
func (m *UserModel) ToEntity() *entity.User {
return &entity.User{
ID:        m.ID,
Email:     m.Email,
Password:  m.Password,
Name:      m.Name,
CreatedAt: m.CreatedAt,
UpdatedAt: m.UpdatedAt,
}
}

// FromEntity converts domain entity to UserModel
func UserModelFromEntity(user *entity.User) *UserModel {
return &UserModel{
ID:        user.ID,
Email:     user.Email,
Password:  user.Password,
Name:      user.Name,
CreatedAt: user.CreatedAt,
UpdatedAt: user.UpdatedAt,
}
}

// TodoModel represents the MongoDB document structure for todos
type TodoModel struct {
ID          primitive.ObjectID `bson:"_id,omitempty"`
UserID      primitive.ObjectID `bson:"user_id"`
Title       string             `bson:"title"`
Description string             `bson:"description"`
Completed   bool               `bson:"completed"`
CreatedAt   time.Time          `bson:"created_at"`
UpdatedAt   time.Time          `bson:"updated_at"`
}

// ToEntity converts TodoModel to domain entity
func (m *TodoModel) ToEntity() *entity.Todo {
return &entity.Todo{
ID:          m.ID,
UserID:      m.UserID,
Title:       m.Title,
Description: m.Description,
Completed:   m.Completed,
CreatedAt:   m.CreatedAt,
UpdatedAt:   m.UpdatedAt,
}
}

// FromEntity converts domain entity to TodoModel
func TodoModelFromEntity(todo *entity.Todo) *TodoModel {
return &TodoModel{
ID:          todo.ID,
UserID:      todo.UserID,
Title:       todo.Title,
Description: todo.Description,
Completed:   todo.Completed,
CreatedAt:   todo.CreatedAt,
UpdatedAt:   todo.UpdatedAt,
}
}
