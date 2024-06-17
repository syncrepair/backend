package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
)

type VehicleUsecase interface {
	Create(ctx context.Context, req VehicleCreateRequest) (string, error)
	Delete(ctx context.Context, id string) error
}

type vehicleUsecase struct {
	repository repository.VehicleRepository
}

func NewVehicleUsecase(repository repository.VehicleRepository) VehicleUsecase {
	return &vehicleUsecase{
		repository: repository,
	}
}

type VehicleCreateRequest struct {
	Make        string
	Model       string
	Year        uint
	VIN         string
	PlateNumber string
	ClientID    string
}

func (uc *vehicleUsecase) Create(ctx context.Context, req VehicleCreateRequest) (string, error) {
	id := util.GenerateID()

	if err := uc.repository.Create(ctx, domain.Vehicle{
		ID:          id,
		Make:        req.Make,
		Model:       req.Model,
		Year:        req.Year,
		VIN:         req.VIN,
		PlateNumber: req.PlateNumber,
		ClientID:    req.ClientID,
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (uc *vehicleUsecase) Delete(ctx context.Context, id string) error {
	return uc.repository.Delete(ctx, id)
}
