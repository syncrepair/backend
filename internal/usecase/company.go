package usecase

import (
	"context"
	"fmt"
)

type CompanyRepository interface {
	Create(ctx context.Context) error
}

type CompanyUsecase struct {
	repository CompanyRepository
}

func NewCompanyUsecase(repository CompanyRepository) *CompanyUsecase {
	return &CompanyUsecase{
		repository: repository,
	}
}

func (u *CompanyUsecase) Create(ctx context.Context) error {
	fmt.Println("CompanyUsecase.Create")

	if err := u.repository.Create(ctx); err != nil {
		return err
	}

	return nil
}
