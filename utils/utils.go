package utils

import (
	"time"

	"github.com/food-order-server/config"
	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateToken(user *models.User) (string, error) {
	config := config.LoadConfig()
	duration, err := time.ParseDuration(config.Token.TokenExpiredTime)
	if err != nil {
		return "", err
	}

	expiredTime := time.Now().Add(duration)
	claims := models.Claims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyPassword(password string, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fiber.ErrBadRequest
	} else {
		return nil
	}
}

func HashPassword(password string) (string, error) {
	config := config.LoadConfig()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), config.EncryptSaltRounds)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func GenerateUuid() string {
	return uuid.NewString()
}
