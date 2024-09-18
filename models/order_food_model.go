package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderFood struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OrderID          string             `bson:"order_id" json:"order_id"`
	FoodID           string             `bson:"food_id" json:"food_id"`
	FoodAmount       int                `bson:"food_amount" json:"food_amount"`
	FoodTotalPrice   float64            `bson:"food_total_price" json:"food_total_price"`
	FoodOptionString string             `bson:"food_option_string" json:"food_option_string"`
	FoodOptionNote   string             `bson:"food_option_note,omitempty" json:"food_option_note,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`

	FoodName          string  `bson:"food_name" json:"food_name"`
	FoodPrice         float64 `bson:"food_price" json:"food_price"`
	FoodPriceDiscount float64 `bson:"food_price_discount" json:"food_price_discount"`
	FoodDescription   string  `bson:"food_description" json:"food_description"`
	FoodImageURL      string  `bson:"food_image_url,omitempty" json:"food_image_url,omitempty"`
	CategoryID        string  `bson:"category_id,omitempty" json:"category_id,omitempty"`
}
