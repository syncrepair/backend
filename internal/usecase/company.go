package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
)

type CompanyUsecase interface {
	Create(ctx context.Context, input CompanyCreateInput) (string, error)
}

type companyUsecase struct {
	repository repository.CompanyRepository
}

func NewCompanyUsecase(repository repository.CompanyRepository) CompanyUsecase {
	return &companyUsecase{
		repository: repository,
	}
}

type CompanyCreateInput struct {
	Name string
}

func (uc *companyUsecase) Create(ctx context.Context, input CompanyCreateInput) (string, error) {
	id := util.GenerateID()

	if err := uc.repository.Create(ctx, domain.Company{
		ID:   id,
		Name: input.Name,
	}); err != nil {
		return "", err
	}

	return id, nil
}
