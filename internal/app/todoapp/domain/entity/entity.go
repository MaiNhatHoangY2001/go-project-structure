package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user entity in the domain
type User struct {
	ID        primitive.ObjectID
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Todo represents a todo item entity in the domain
type Todo struct {
	ID          primitive.ObjectID
	UserID      primitive.ObjectID
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
