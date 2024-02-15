package usecase

import (
	"context"
	"fmt"
)

type CompanyUsecase struct{}

func NewCompanyUsecase() *CompanyUsecase {
	return &CompanyUsecase{}
}

func (u *CompanyUsecase) Create(ctx context.Context) error {
	fmt.Println("CompanyUsecase.Create")

	return nil
}
