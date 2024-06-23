package domain

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client already exists")
)

type Client struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	Vehicles    []ClientVehicle    `json:"vehicles" bson:"vehicles"`
	Settings    ClientSettings     `json:"settings" bson:"settings"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id"`
}

type ClientVehicle struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Make        string             `json:"make" bson:"make"`
	Model       string             `json:"model" bson:"model"`
	Year        uint               `json:"year" bson:"year"`
	VIN         string             `json:"vin" bson:"vin"`
	Distance    uint               `json:"distance" bson:"distance"`
	PlateNumber string             `json:"plate_number" bson:"plate_number"`
}

type ClientSettings struct {
	Discount int `json:"discount" bson:"discount"`
}
