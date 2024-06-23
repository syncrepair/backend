package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VehicleUsecase struct {
	repository repository.ClientRepository
}

func NewVehicleUsecase(repository repository.ClientRepository) *VehicleUsecase {
	return &VehicleUsecase{
		repository: repository,
	}
}

type VehicleCreateInput struct {
	Make        string
	Model       string
	Year        uint
	VIN         string
	Distance    uint
	PlateNumber string
}

func (uc *VehicleUsecase) Create(ctx context.Context, clientID string, input VehicleCreateInput) (string, error) {
	clientId, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return "", err
	}

	vehicleId := primitive.NewObjectID()

	if err := uc.repository.AddVehicle(ctx, clientId, domain.ClientVehicle{
		ID:          vehicleId,
		Make:        input.Make,
		Model:       input.Model,
		Year:        input.Year,
		VIN:         input.VIN,
		Distance:    input.Distance,
		PlateNumber: input.PlateNumber,
	}); err != nil {
		return "", err
	}

	return vehicleId.Hex(), nil
}

func (uc *VehicleUsecase) GetAll(ctx context.Context, clientID string) ([]domain.ClientVehicle, error) {
	clientId, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, err
	}

	vehicles, err := uc.repository.GetAllVehicles(ctx, clientId)
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (uc *VehicleUsecase) GetByID(ctx context.Context, clientID, vehicleID string) (domain.ClientVehicle, error) {
	clientId, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return domain.ClientVehicle{}, err
	}

	vehicleId, err := primitive.ObjectIDFromHex(vehicleID)
	if err != nil {
		return domain.ClientVehicle{}, err
	}

	vehicle, err := uc.repository.GetVehicleByID(ctx, clientId, vehicleId)
	if err != nil {
		return domain.ClientVehicle{}, err
	}

	return vehicle, nil
}

type VehicleUpdateInput struct {
	Make        string
	Model       string
	Year        uint
	VIN         string
	Distance    uint
	PlateNumber string
}

func (uc *VehicleUsecase) Update(ctx context.Context, clientID, vehicleID string, input VehicleUpdateInput) error {
	clientId, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return err
	}

	vehicleId, err := primitive.ObjectIDFromHex(vehicleID)
	if err != nil {
		return err
	}

	if err := uc.repository.UpdateVehicle(ctx, clientId, domain.ClientVehicle{
		ID:          vehicleId,
		Make:        input.Make,
		Model:       input.Model,
		Year:        input.Year,
		VIN:         input.VIN,
		Distance:    input.Distance,
		PlateNumber: input.PlateNumber,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *VehicleUsecase) Delete(ctx context.Context, clientID, vehicleID string) error {
	clientId, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return err
	}

	vehicleId, err := primitive.ObjectIDFromHex(vehicleID)
	if err != nil {
		return err
	}

	if err := uc.repository.DeleteVehicle(ctx, clientId, vehicleId); err != nil {
		return err
	}

	return nil
}
