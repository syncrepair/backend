package model

import (
	"errors"
	"time"
)

var (
	ErrCompanyNameLength = errors.New("name must be at least 2 characters long")
)

type Company struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompanyCreateInput struct {
	Name string `json:"name"`
}

func (c *CompanyCreateInput) Validate() error {
	if len(c.Name) < 2 {
		return ErrCompanyNameLength
	}

	return nil
}
