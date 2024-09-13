package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LatLng struct {
	Lat float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	Lng float64 `bson:"lng,omitempty" json:"lng,omitempty"`
}

type Location struct {
	Address string `bson:"address,omitempty" json:"address,omitempty"`
	LatLng  LatLng `bson:"latlng,omitempty" json:"latlng,omitempty"`
}

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID          string             `bson:"user_id" json:"user_id"`
	Username        string             `bson:"username" json:"username"`
	Password        string             `bson:"password" json:"password"`
	Name            string             `bson:"name" json:"name"`
	Role            string             `bson:"role" json:"role"`
	Balance         float64            `bson:"balance" json:"balance"`
	Location        Location           `bson:"location,omitempty" json:"location,omitempty"`
	ProfileImageURL string             `bson:"profile_image_url,omitempty" json:"profile_image_url,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
