package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	UserCollectionName = "users"

	UserWorkerRole Role = "worker"
	UserOwnerRole  Role = "owner"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	Role        Role               `json:"role" bson:"role"`
	IsConfirmed bool               `json:"isConfirmed" bson:"is_confirmed"`
}

type Role string

type UserTokens struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
