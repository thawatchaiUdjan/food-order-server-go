package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	OrderID             string         `bson:"order_id" json:"order_id"`
	UserID              string         `bson:"user_id" json:"user_id"`
	SubtotalPrice       float64        `bson:"subtotal_price" json:"subtotal_price"`
	TotalPrice          float64        `bson:"total_price" json:"total_price"`
	OrderAddress        string         `bson:"order_address" json:"order_address"`
	OrderStatus         OrderStatus    `bson:"order_status" json:"order_status"`
	OrderDeliveryOption DeliveryOption `bson:"order_delivery_option" json:"order_delivery_option"`
	CreatedAt           time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time      `bson:"updated_at" json:"updated_at"`
}

type OrderCreate struct {
	OrderID             string             `bson:"order_id" json:"order_id"`
	UserID              string             `bson:"user_id" json:"user_id"`
	SubtotalPrice       float64            `bson:"subtotal_price" json:"subtotal_price"`
	TotalPrice          float64            `bson:"total_price" json:"total_price"`
	OrderAddress        string             `bson:"order_address" json:"order_address"`
	OrderStatus         primitive.ObjectID `bson:"order_status" json:"order_status"`
	OrderDeliveryOption primitive.ObjectID `bson:"order_delivery_option" json:"order_delivery_option"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}

type OrderAll struct {
	OrderID         string             `bson:"order_id" json:"order_id"`
	SubtotalPrice   float64            `bson:"subtotal_price" json:"subtotal_price"`
	TotalPrice      float64            `bson:"total_price" json:"total_price"`
	OrderAddress    string             `bson:"order_address" json:"order_address"`
	OrderStatus     primitive.ObjectID `bson:"order_status" json:"order_status"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	UserID          string             `bson:"user_id" json:"user_id"`
	Username        string             `bson:"username" json:"username"`
	Password        string             `bson:"password" json:"password"`
	Name            string             `bson:"name" json:"name"`
	Role            string             `bson:"role" json:"role"`
	Balance         float64            `bson:"balance" json:"balance"`
	Location        Location           `bson:"location,omitempty" json:"location,omitempty"`
	ProfileImageURL string             `bson:"profile_image_url,omitempty" json:"profile_image_url,omitempty"`
	StatusName      string             `bson:"status_name" json:"status_name"`
	StatusValue     int                `bson:"status_value" json:"status_value"`
	DeliveryOption  DeliveryOption     `bson:"delivery_option" json:"delivery_option"`
	Foods           []OrderFood        `bson:"foods" json:"foods"`
}
