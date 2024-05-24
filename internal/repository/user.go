package repository

import (
	"context"
	"fmt"
	"github.com/syncrepair/backend/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type User interface {
	Create(context.Context, *entity.User) (*entity.User, error)
}

type user struct {
	collection *mongo.Collection
}

func NewUser(collection *mongo.Collection) User {
	return &user{
		collection: collection,
	}
}

func (u *user) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	result, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return user, nil
}
