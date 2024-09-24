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
	FoodOptionString []string           `bson:"food_option_string" json:"food_option_string"`
	FoodOptionNote   string             `bson:"food_option_note,omitempty" json:"food_option_note,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`

	FoodName          string  `bson:"food_name,omitempty" json:"food_name,omitempty"`
	FoodPrice         float64 `bson:"food_price,omitempty" json:"food_price,omitempty"`
	FoodPriceDiscount float64 `bson:"food_price_discount,omitempty" json:"food_price_discount,omitempty"`
	FoodDescription   string  `bson:"food_description,omitempty" json:"food_description,omitempty"`
	FoodImageURL      string  `bson:"food_image_url,omitempty" json:"food_image_url,omitempty"`
	CategoryID        string  `bson:"category_id,omitempty" json:"category_id,omitempty"`
}

type OrderFoodCreate struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OrderID          string             `bson:"order_id" json:"order_id"`
	FoodID           string             `bson:"food_id" json:"food_id"`
	FoodAmount       int                `bson:"food_amount" json:"food_amount"`
	FoodTotalPrice   float64            `bson:"food_total_price" json:"food_total_price"`
	FoodOptionString []string           `bson:"food_option_string" json:"food_option_string"`
	FoodOptionNote   string             `bson:"food_option_note,omitempty" json:"food_option_note,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
