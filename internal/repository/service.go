package repository

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceRepository interface {
	Create(ctx context.Context, service domain.Service) error
	GetAll(ctx context.Context, companyID primitive.ObjectID) ([]domain.Service, error)
	GetByID(ctx context.Context, id, companyID primitive.ObjectID) (domain.Service, error)
	Update(ctx context.Context, service domain.Service) error
	Delete(ctx context.Context, id, companyID primitive.ObjectID) error
}

type serviceRepository struct {
	db *mongo.Collection
}

func NewServiceRepository(db *mongo.Database) ServiceRepository {
	return &serviceRepository{
		db: db.Collection(servicesCollectionName),
	}
}

func (r *serviceRepository) Create(ctx context.Context, service domain.Service) error {
	_, err := r.db.InsertOne(ctx, service)
	if err != nil {
		return err
	}

	return nil
}

func (r *serviceRepository) GetAll(ctx context.Context, companyID primitive.ObjectID) ([]domain.Service, error) {
	cur, err := r.db.Find(ctx, bson.M{"company_id": companyID})
	if err != nil {
		return nil, err
	}

	var services []domain.Service
	if err := cur.All(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceRepository) GetByID(ctx context.Context, id, companyID primitive.ObjectID) (domain.Service, error) {
	var service domain.Service
	if err := r.db.FindOne(ctx, bson.M{"_id": id, "company_id": companyID}).Decode(&service); err != nil {
		if mongodb.IsNotFound(err) {
			return domain.Service{}, domain.ErrServiceNotFound
		}

		return domain.Service{}, err
	}

	return service, nil
}

func (r *serviceRepository) Update(ctx context.Context, service domain.Service) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": service.ID}, bson.M{"$set": service})
	if err != nil {
		return err
	}

	return nil
}

func (r *serviceRepository) Delete(ctx context.Context, id, companyID primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id, "company_id": companyID})

	return err
}
