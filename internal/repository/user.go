package repository

import (
	"context"
	"fmt"
	"github.com/syncrepair/backend/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(context.Context, *entity.User) (*entity.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return user, nil
}
