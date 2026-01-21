package repository

import (
	"context"
	"time"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	model := UserModelFromEntity(user)
	result, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model UserModel
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&model)
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *userRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	var model UserModel
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}
