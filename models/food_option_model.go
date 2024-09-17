package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodOption struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OptionName        string             `bson:"option_name" json:"option_name"`
	OptionDescription string             `bson:"option_description" json:"option_description"`
	OptionType        string             `bson:"option_type" json:"option_type"`
	OptionChoices     []OptionChoice     `bson:"option_choices" json:"option_choices"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}
