package postgres

import (
	"context"
	"fmt"
	"github.com/syncrepair/backend/internal/model"
)

type CompanyRepository struct{}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{}
}

func (r *CompanyRepository) Create(ctx context.Context, company *model.Company) error {
	fmt.Println(company)

	return nil
}
