package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
)

type ServiceUsecase interface {
	Create(ctx context.Context, req ServiceCreateRequest) (string, error)
	Delete(ctx context.Context, id string) error
}

type serviceUsecase struct {
	repository repository.ServiceRepository
}

func NewServiceUsecase(repository repository.ServiceRepository) ServiceUsecase {
	return &serviceUsecase{
		repository: repository,
	}
}

type ServiceCreateRequest struct {
	Name        string
	Description string
	Price       float64
	CompanyID   string
}

func (uc *serviceUsecase) Create(ctx context.Context, req ServiceCreateRequest) (string, error) {
	id := util.GenerateID()

	if err := uc.repository.Create(ctx, domain.Service{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CompanyID:   req.CompanyID,
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (uc *serviceUsecase) Delete(ctx context.Context, id string) error {
	return uc.repository.Delete(ctx, id)
}
