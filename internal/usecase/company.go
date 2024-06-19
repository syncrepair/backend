package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
	"time"
)

type CompanyUsecase interface {
	Create(ctx context.Context, req CompanyCreateRequest) (string, error)
	Delete(ctx context.Context, id string) error
}

type companyUsecase struct {
	repository repository.CompanyRepository
}

func NewCompanyUsecase(repository repository.CompanyRepository) CompanyUsecase {
	return &companyUsecase{
		repository: repository,
	}
}

type CompanyCreateRequest struct {
	Name      string
	OpenTime  time.Time
	CloseTime time.Time
}

func (uc *companyUsecase) Create(ctx context.Context, req CompanyCreateRequest) (string, error) {
	id := util.GenerateID()

	if err := uc.repository.Create(ctx, domain.Company{
		ID:        id,
		Name:      req.Name,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (uc *companyUsecase) Delete(ctx context.Context, id string) error {
	return uc.repository.Delete(ctx, id)
}
