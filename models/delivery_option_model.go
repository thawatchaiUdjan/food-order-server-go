package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryOption struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DeliveryName        string             `bson:"delivery_name" json:"delivery_name"`
	DeliveryDescription string             `bson:"delivery_description" json:"delivery_description"`
	DeliveryCost        float64            `bson:"delivery_cost" json:"delivery_cost"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}
