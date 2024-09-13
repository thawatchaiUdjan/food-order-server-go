package models

type MessageRes struct {
	Message string `json:"message"`
}

type UserDataRes struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
