package repository

import (
	"context"
	"fmt"
	"github.com/syncrepair/backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(context.Context, *domain.User) (*domain.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection(domain.UserCollectionName),
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return user, nil
}
