package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	StatusName  string             `bson:"status_name" json:"status_name"`
	StatusValue int                `bson:"status_value" json:"status_value"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
