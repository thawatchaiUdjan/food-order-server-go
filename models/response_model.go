package models

type MessageRes struct {
	Message string `json:"message"`
}

type UserDataRes struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type UserUpdateRes struct {
	User    User   `json:"user"`
	Message string `json:"message"`
}

type FoodDataRes struct {
	Food    Food   `json:"food"`
	Message string `json:"message"`
}

type FoodOrderRes struct {
	Order Order       `json:"order"`
	Foods []OrderFood `json:"foods"`
}

type OrderDataRes struct {
	FoodOrder FoodOrderRes `json:"foodOrder"`
	User      UserDataRes  `json:"user"`
	Message   string       `json:"message"`
}
