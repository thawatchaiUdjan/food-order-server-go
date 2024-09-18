package models

type UserLoginReq struct {
	Username string `bson:"username" json:"username" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
}

type UserRegisterReq struct {
	Username string `bson:"username" json:"username" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
	Name     string `bson:"name" json:"name" validate:"required"`
}

type FoodCartOption struct {
	OptionString string `json:"option_string" bson:"option_string"`
	OptionNote   string `json:"option_note" bson:"option_note"`
}

type FoodCart struct {
	Food   Food           `json:"food" bson:"food" validate:"required"`
	Amount int            `json:"amount" bson:"amount" validate:"required"`
	Total  float64        `json:"total" bson:"total" validate:"required"`
	Option FoodCartOption `json:"option" bson:"option"`
}

type OrderReq struct {
	Foods []FoodCart  `json:"foods" bson:"foods" validate:"required"`
	Order CreateOrder `json:"order" bson:"order" validate:"required"`
}
