package domain

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrServiceNotFound = errors.New("service not found")
)

type Service struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Duration    uint               `json:"duration" bson:"duration"`
	Price       float64            `json:"price" bson:"price"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id"`
}
