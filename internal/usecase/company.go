package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/model"
	"github.com/syncrepair/backend/internal/utils/id"
	"time"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *model.Company) error
}

type CompanyUsecase struct {
	repository CompanyRepository
}

func NewCompanyUsecase(repository CompanyRepository) *CompanyUsecase {
	return &CompanyUsecase{
		repository: repository,
	}
}

func (u *CompanyUsecase) Create(ctx context.Context, input *model.CompanyCreateInput) error {
	company := &model.Company{
		ID:        id.Generate(),
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.repository.Create(ctx, company); err != nil {
		return err
	}

	return nil
}
