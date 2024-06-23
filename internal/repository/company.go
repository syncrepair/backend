package repository

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyRepository interface {
	Create(ctx context.Context, company domain.Company) error
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.Company, error)
	Update(ctx context.Context, company domain.Company) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type companyRepository struct {
	db *mongo.Collection
}

func NewCompanyRepository(db *mongo.Database) CompanyRepository {
	return &companyRepository{
		db: db.Collection(companiesCollectionName),
	}
}

func (r *companyRepository) Create(ctx context.Context, company domain.Company) error {
	_, err := r.db.InsertOne(ctx, company)
	if err != nil {
		return err
	}

	return nil
}

func (r *companyRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Company, error) {
	var company domain.Company
	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&company); err != nil {
		if mongodb.IsNotFound(err) {
			return domain.Company{}, domain.ErrCompanyNotFound
		}

		return domain.Company{}, err
	}

	return company, nil
}

func (r *companyRepository) Update(ctx context.Context, company domain.Company) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": company.ID}, bson.M{"$set": company})
	if err != nil {
		return err
	}

	return nil
}

func (r *companyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})

	return err
}
