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

type UserGoogleLoginReq struct {
	Code string `bson:"code" json:"code" validate:"required"`
}

type FoodCartOption struct {
	OptionString []string `json:"option_string" bson:"option_string"`
	OptionNote   string   `json:"option_note" bson:"option_note"`
}

type FoodCart struct {
	Food   Food           `json:"food" bson:"food"`
	Amount int            `json:"amount" bson:"amount"`
	Total  float64        `json:"total" bson:"total"`
	Option FoodCartOption `json:"option" bson:"option"`
}

type OrderReq struct {
	Foods []FoodCart  `json:"foods" bson:"foods" validate:"required"`
	Order OrderCreate `json:"order" bson:"order" validate:"required"`
}

type FoodReq struct {
	FoodName          string   `bson:"food_name" json:"food_name" form:"food_name" validate:"required"`
	FoodPrice         float64  `bson:"food_price" json:"food_price" form:"food_price" validate:"required"`
	FoodPriceDiscount float64  `bson:"food_price_discount" json:"food_price_discount" form:"food_price_discount" validate:"required"`
	FoodDescription   string   `bson:"food_description" json:"food_description" form:"food_description" validate:"required"`
	FoodImageURL      string   `bson:"food_image_url" json:"food_image_url"`
	CategoryID        string   `bson:"category_id" json:"category_id" form:"category_id" validate:"required"`
	FoodOptions       []string `bson:"food_options" json:"food_options" form:"food_options[]" validate:"required"`
}
