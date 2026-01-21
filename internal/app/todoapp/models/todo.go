package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represents a todo item
type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"userId"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// TodoCreateRequest represents a request to create a todo
type TodoCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// TodoUpdateRequest represents a request to update a todo
type TodoUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

// TodoListQuery represents query parameters for listing todos
type TodoListQuery struct {
	Page      int  `form:"page" binding:"min=1"`
	PageSize  int  `form:"pageSize" binding:"min=1,max=100"`
	Completed *bool `form:"completed"`
}
