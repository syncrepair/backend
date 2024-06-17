package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
)

type ClientUsecase interface {
	Create(ctx context.Context, req ClientCreateRequest) (string, error)
	Delete(ctx context.Context, id string) error
}

type clientUsecase struct {
	repository repository.ClientRepository
}

func NewClientUsecase(repository repository.ClientRepository) ClientUsecase {
	return &clientUsecase{
		repository: repository,
	}
}

type ClientCreateRequest struct {
	Name        string
	PhoneNumber string
	CompanyID   string
}

func (uc *clientUsecase) Create(ctx context.Context, req ClientCreateRequest) (string, error) {
	id := util.GenerateID()

	if err := uc.repository.Create(ctx, domain.Client{
		ID:          id,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		CompanyID:   req.CompanyID,
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (uc *clientUsecase) Delete(ctx context.Context, id string) error {
	return uc.repository.Delete(ctx, id)
}
