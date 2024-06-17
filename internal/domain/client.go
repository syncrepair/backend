package domain

import "errors"

var (
	ErrClientNotFound = errors.New("клієнта не знайдено")
)

type Client struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	CompanyID   string `json:"company_id"`
}
