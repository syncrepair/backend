package repository

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepository interface {
	Create(ctx context.Context, client domain.Client) error
	GetAll(ctx context.Context, companyID primitive.ObjectID) ([]domain.Client, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.Client, error)
	Update(ctx context.Context, client domain.Client) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	AddVehicle(ctx context.Context, clientID primitive.ObjectID, vehicle domain.ClientVehicle) error
	GetAllVehicles(ctx context.Context, clientID primitive.ObjectID) ([]domain.ClientVehicle, error)
	GetVehicleByID(ctx context.Context, clientID, vehicleID primitive.ObjectID) (domain.ClientVehicle, error)
	UpdateVehicle(ctx context.Context, clientID primitive.ObjectID, vehicle domain.ClientVehicle) error
	DeleteVehicle(ctx context.Context, clientID, vehicleID primitive.ObjectID) error
}

type clientRepository struct {
	db *mongo.Collection
}

func NewClientRepository(db *mongo.Database) ClientRepository {
	return &clientRepository{
		db: db.Collection(clientsCollectionName),
	}
}

func (r *clientRepository) Create(ctx context.Context, client domain.Client) error {
	_, err := r.db.InsertOne(ctx, client)
	if err != nil {
		if mongodb.IsDuplicate(err) {
			return domain.ErrClientAlreadyExists
		}

		return err
	}

	return nil
}

func (r *clientRepository) GetAll(ctx context.Context, companyID primitive.ObjectID) ([]domain.Client, error) {
	cur, err := r.db.Find(ctx, bson.M{"company_id": companyID})
	if err != nil {
		return nil, err
	}

	var clients []domain.Client
	if err := cur.All(ctx, &clients); err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *clientRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Client, error) {
	var client domain.Client
	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&client); err != nil {
		if mongodb.IsNotFound(err) {
			return domain.Client{}, domain.ErrClientNotFound
		}

		return domain.Client{}, err
	}

	return client, nil
}

func (r *clientRepository) Update(ctx context.Context, client domain.Client) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": client.ID}, bson.M{"$set": client})
	if err != nil {
		return err
	}

	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *clientRepository) AddVehicle(ctx context.Context, clientID primitive.ObjectID, vehicle domain.ClientVehicle) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": clientID}, bson.M{"$push": bson.M{"vehicles": vehicle}})

	return err
}

func (r *clientRepository) GetAllVehicles(ctx context.Context, clientID primitive.ObjectID) ([]domain.ClientVehicle, error) {
	client, err := r.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return client.Vehicles, nil
}

func (r *clientRepository) GetVehicleByID(ctx context.Context, clientID, vehicleID primitive.ObjectID) (domain.ClientVehicle, error) {
	var client domain.Client
	if err := r.db.FindOne(ctx, bson.M{"_id": clientID}).Decode(&client); err != nil {
		if mongodb.IsNotFound(err) {
			return domain.ClientVehicle{}, domain.ErrClientNotFound
		}

		return domain.ClientVehicle{}, err
	}

	return client.Vehicles[0], nil
}

func (r *clientRepository) UpdateVehicle(ctx context.Context, clientID primitive.ObjectID, vehicle domain.ClientVehicle) error {
	_, err := r.db.UpdateOne(ctx,
		bson.M{
			"_id":          clientID,
			"vehicles._id": vehicle.ID,
		},
		bson.M{"$set": bson.M{
			"vehicles.$.make":         vehicle.Make,
			"vehicles.$.model":        vehicle.Model,
			"vehicles.$.year":         vehicle.Year,
			"vehicles.$.vin":          vehicle.VIN,
			"vehicles.$.distance":     vehicle.Distance,
			"vehicles.$.plate_number": vehicle.PlateNumber,
		}},
	)

	return err
}

func (r *clientRepository) DeleteVehicle(ctx context.Context, clientID, vehicleID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": clientID},
		bson.M{"$pull": bson.M{"vehicles": bson.M{"_id": vehicleID}}},
	)

	return err
}
