package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FoodID            string             `bson:"food_id" json:"food_id"`
	FoodName          string             `bson:"food_name" json:"food_name"`
	FoodPrice         float64            `bson:"food_price" json:"food_price"`
	FoodPriceDiscount float64            `bson:"food_price_discount" json:"food_price_discount"`
	FoodDescription   string             `bson:"food_description" json:"food_description"`
	FoodImageURL      string             `bson:"food_image_url" json:"food_image_url"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	CategoryID        string             `bson:"category_id,omitempty" json:"category_id,omitempty"`
	FoodOptions       []FoodOption       `bson:"food_options,omitempty" json:"food_options,omitempty"`
}

type FoodCreate struct {
	FoodID            string               `bson:"food_id" json:"food_id"`
	FoodName          string               `bson:"food_name" json:"food_name"`
	FoodPrice         float64              `bson:"food_price" json:"food_price"`
	FoodPriceDiscount float64              `bson:"food_price_discount" json:"food_price_discount"`
	FoodDescription   string               `bson:"food_description" json:"food_description"`
	FoodImageURL      string               `bson:"food_image_url" json:"food_image_url"`
	CreatedAt         time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time            `bson:"updated_at" json:"updated_at"`
	CategoryID        string               `bson:"category_id" json:"category_id"`
	FoodOptions       []primitive.ObjectID `bson:"food_options" json:"food_options"`
}
