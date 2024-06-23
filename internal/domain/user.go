package domain

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user with such credentials was not found")
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	IsConfirmed bool               `json:"is_confirmed" bson:"is_confirmed"`
	Session     UserSession        `json:"session" bson:"session"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id"`
}

type UserSession struct {
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" bson:"expires_at"`
}
