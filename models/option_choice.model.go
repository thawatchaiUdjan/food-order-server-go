package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OptionChoice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ChoiceName  string             `bson:"choice_name" json:"choice_name"`
	ChoicePrice float64            `bson:"choice_price" json:"choice_price"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
