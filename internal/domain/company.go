package domain

import "errors"

var (
	ErrCompanyNotFound = errors.New("компанію не знайдено")
)

type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
