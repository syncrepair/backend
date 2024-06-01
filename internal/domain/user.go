package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("користувача не знайдено")
	ErrUserAlreadyExists = errors.New("користувач вже існує")
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsConfirmed bool   `json:"is_confirmed"`
}

type UserTokens struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
