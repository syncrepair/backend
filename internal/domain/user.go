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
	CompanyID   string `json:"company_id"`
	IsConfirmed bool   `json:"is_confirmed"`
}
