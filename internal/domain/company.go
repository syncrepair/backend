package domain

import (
	"errors"
	"time"
)

var (
	ErrCompanyNotFound = errors.New("компанію не знайдено")
)

type Company struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}
