package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceUsecase struct {
	repository repository.ServiceRepository
}

func NewServiceUsecase(repository repository.ServiceRepository) *ServiceUsecase {
	return &ServiceUsecase{
		repository: repository,
	}
}

type ServiceCreateInput struct {
	Name        string
	Description string
	Duration    uint
	Price       float64
	CompanyID   string
}

func (uc *ServiceUsecase) Create(ctx context.Context, input ServiceCreateInput) (string, error) {
	id := primitive.NewObjectID()

	companyID, err := primitive.ObjectIDFromHex(input.CompanyID)
	if err != nil {
		return "", err
	}

	if err := uc.repository.Create(ctx, domain.Service{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Duration:    input.Duration,
		Price:       input.Price,
		CompanyID:   companyID,
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (uc *ServiceUsecase) GetAll(ctx context.Context, companyID string) ([]domain.Service, error) {
	id, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, err
	}

	services, err := uc.repository.GetAll(ctx, id)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (uc *ServiceUsecase) GetByID(ctx context.Context, id, companyID string) (domain.Service, error) {
	serviceId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Service{}, err
	}

	companyId, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return domain.Service{}, err
	}

	return uc.repository.GetByID(ctx, serviceId, companyId)
}

type ServiceUpdateInput struct {
	Name        string
	Description string
	Duration    uint
	Price       float64
	CompanyID   string
}

func (uc *ServiceUsecase) Update(ctx context.Context, id string, input ServiceUpdateInput) error {
	serviceId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	companyId, err := primitive.ObjectIDFromHex(input.CompanyID)
	if err != nil {
		return err
	}

	return uc.repository.Update(ctx, domain.Service{
		ID:          serviceId,
		Name:        input.Name,
		Description: input.Description,
		Duration:    input.Duration,
		Price:       input.Price,
		CompanyID:   companyId,
	})
}

func (uc *ServiceUsecase) Delete(ctx context.Context, id, companyID string) error {
	serviceId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	companyId, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return err
	}

	return uc.repository.Delete(ctx, serviceId, companyId)
}
