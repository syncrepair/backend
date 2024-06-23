package repository

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	GetByCredentials(ctx context.Context, email string, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	Confirm(ctx context.Context, id primitive.ObjectID) error
	SetSession(ctx context.Context, id primitive.ObjectID, session domain.UserSession) error
}

type userRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db.Collection(usersCollectionName),
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		if mongodb.IsDuplicate(err) {
			return domain.ErrUserAlreadyExists
		}

		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if mongodb.IsNotFound(err) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByCredentials(ctx context.Context, email string, password string) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if mongodb.IsNotFound(err) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{
		"session.refresh_token": refreshToken,
		"session.expires_at":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if mongodb.IsNotFound(err) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user domain.User) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}})
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *userRepository) Confirm(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_confirmed": true}})

	return err
}

func (r *userRepository) SetSession(ctx context.Context, id primitive.ObjectID, session domain.UserSession) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"session": session}})

	return err
}
