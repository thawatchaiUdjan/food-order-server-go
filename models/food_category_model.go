package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodCategory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CategoryID   string             `bson:"category_id" json:"category_id"`
	CategoryName string             `bson:"category_name" json:"category_name"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
