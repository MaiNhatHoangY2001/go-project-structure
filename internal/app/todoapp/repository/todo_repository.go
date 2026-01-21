package repository

import (
	"context"
	"time"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type todoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(db *mongo.Database) domain.TodoRepository {
	return &todoRepository{
		collection: db.Collection("todos"),
	}
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	model := TodoModelFromEntity(todo)
	result, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}
	todo.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *todoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Todo, error) {
	var model TodoModel
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *todoRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int, completed *bool) ([]*entity.Todo, int64, error) {
	filter := bson.M{"user_id": userID}
	if completed != nil {
		filter["completed"] = *completed
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var models []*TodoModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, 0, err
	}

	todos := make([]*entity.Todo, len(models))
	for i, model := range models {
		todos[i] = model.ToEntity()
	}

	return todos, total, nil
}

func (r *todoRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *todoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
