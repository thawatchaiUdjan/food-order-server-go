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

type CreateOrder struct {
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
