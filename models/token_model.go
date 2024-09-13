package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}

type UserReq struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
