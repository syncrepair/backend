package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyUsecase struct {
	repository repository.CompanyRepository
}

func NewCompanyUsecase(repository repository.CompanyRepository) *CompanyUsecase {
	return &CompanyUsecase{
		repository: repository,
	}
}

type CompanyCreateInput struct {
	Name     string
	Logo     string
	Location domain.CompanyLocation
	Settings domain.CompanySettings
}

func (uc *CompanyUsecase) Create(ctx context.Context, input CompanyCreateInput) (string, error) {
	id := primitive.NewObjectID()

	if err := uc.repository.Create(ctx, domain.Company{
		ID:       id,
		Name:     input.Name,
		Logo:     input.Logo,
		Location: input.Location,
		Settings: input.Settings,
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (uc *CompanyUsecase) GetByID(ctx context.Context, id string) (domain.Company, error) {
	companyID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Company{}, err
	}

	return uc.repository.GetByID(ctx, companyID)
}

type CompanyUpdateInput struct {
	Name     string
	Logo     string
	Location domain.CompanyLocation
	Settings domain.CompanySettings
}

func (uc *CompanyUsecase) Update(ctx context.Context, id string, input CompanyUpdateInput) error {
	companyID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return uc.repository.Update(ctx, domain.Company{
		ID:       companyID,
		Name:     input.Name,
		Logo:     input.Logo,
		Location: input.Location,
		Settings: input.Settings,
	})
}

func (uc *CompanyUsecase) Delete(ctx context.Context, id string) error {
	companyID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return uc.repository.Delete(ctx, companyID)
}
