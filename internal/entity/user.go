package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	IsConfirmed bool               `bson:"is_confirmed"`
}

type UserTokens struct {
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
}
