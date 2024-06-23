package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientUsecase struct {
	repository repository.ClientRepository
}

func NewClientUsecase(repository repository.ClientRepository) *ClientUsecase {
	return &ClientUsecase{
		repository: repository,
	}
}

type ClientCreateInput struct {
	Name        string
	PhoneNumber string
	Vehicles    []domain.ClientVehicle
	Settings    domain.ClientSettings
	CompanyID   string
}

func (uc *ClientUsecase) Create(ctx context.Context, input ClientCreateInput) (string, error) {
	id := primitive.NewObjectID()

	companyID, err := primitive.ObjectIDFromHex(input.CompanyID)
	if err != nil {
		return "", err
	}

	if err := uc.repository.Create(ctx, domain.Client{
		ID:          id,
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Vehicles:    input.Vehicles,
		Settings:    input.Settings,
		CompanyID:   companyID,
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (uc *ClientUsecase) GetAll(ctx context.Context, companyID string) ([]domain.Client, error) {
	id, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, err
	}

	return uc.repository.GetAll(ctx, id)
}

func (uc *ClientUsecase) GetByID(ctx context.Context, id string) (domain.Client, error) {
	clientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Client{}, err
	}

	return uc.repository.GetByID(ctx, clientID)
}

type ClientUpdateInput struct {
	Name        string
	PhoneNumber string
	Vehicles    []domain.ClientVehicle
	Settings    domain.ClientSettings
	CompanyID   string
}

func (uc *ClientUsecase) Update(ctx context.Context, id string, input ClientUpdateInput) error {
	clientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	companyID, err := primitive.ObjectIDFromHex(input.CompanyID)
	if err != nil {
		return err
	}

	return uc.repository.Update(ctx, domain.Client{
		ID:          clientID,
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Vehicles:    input.Vehicles,
		Settings:    input.Settings,
		CompanyID:   companyID,
	})
}

func (uc *ClientUsecase) Delete(ctx context.Context, id string) error {
	clientID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return uc.repository.Delete(ctx, clientID)
}
