package model

import (
	"errors"
	"net/mail"
	"time"
)

var (
	ErrUserNameLength     = errors.New("name must be at least 2 characters long")
	ErrUserPasswordLength = errors.New("password must be at least 8 characters long")
	ErrUserEmailInvalid   = errors.New("invalid email address")
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CompanyID string    `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSignUpInput struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyCode string `json:"company_code"`
}

func (u *UserSignUpInput) Validate() error {
	if len(u.Name) < 2 {
		return ErrUserNameLength
	}

	if len(u.Password) < 8 {
		return ErrUserPasswordLength
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return ErrUserEmailInvalid
	}

	return nil
}
