package postgres

import (
	"context"
	"fmt"
)

type CompanyRepository struct{}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{}
}

func (r *CompanyRepository) Create(ctx context.Context) error {
	fmt.Println("CompanyRepository.Create")

	return nil
}
