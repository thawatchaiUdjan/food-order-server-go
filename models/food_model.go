package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FoodID            string             `bson:"food_id" json:"food_id"`
	FoodName          string             `bson:"food_name" json:"food_name" validate:"required"`
	FoodPrice         float64            `bson:"food_price" json:"food_price" validate:"required"`
	FoodPriceDiscount float64            `bson:"food_price_discount" json:"food_price_discount" validate:"required"`
	FoodDescription   string             `bson:"food_description" json:"food_description" validate:"required"`
	FoodImageURL      string             `bson:"food_image_url,omitempty" json:"food_image_url,omitempty"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	CategoryID        string             `bson:"category_id,omitempty" json:"category_id,omitempty"`
	FoodOptions       []FoodOption       `bson:"food_options,omitempty" json:"food_options,omitempty"`
}
