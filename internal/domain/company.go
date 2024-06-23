package domain

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	ErrCompanyNotFound = errors.New("company not found")
)

type Company struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Logo     string             `json:"logo" bson:"logo"`
	Location CompanyLocation    `json:"location" bson:"location"`
	Settings CompanySettings    `json:"settings" bson:"settings"`
}

type CompanySettings struct {
	OpenTime        time.Time       `json:"open_time" bson:"open_time"`
	CloseTime       time.Time       `json:"close_time" bson:"close_time"`
	Currency        string          `json:"currency" bson:"currency"`
	MeasurementUnit MeasurementUnit `json:"measurement_unit" bson:"measurement_unit"`
}

type MeasurementUnit string

const (
	KilometresMeasurementUnit MeasurementUnit = "km"
	MilesMeasurementUnit      MeasurementUnit = "mi"
)

type CompanyLocation struct {
	Longitude float64 `json:"longitude" bson:"longitude"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
}
