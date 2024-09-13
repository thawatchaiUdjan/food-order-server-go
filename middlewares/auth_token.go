package middlewares

import (
	"strings"

	"github.com/food-order-server/config"
	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthToken(c fiber.Ctx) error {
	config := config.LoadConfig()
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No token provided. Access denied")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err == jwt.ErrTokenExpired {
		return fiber.NewError(fiber.StatusUnauthorized, "Token has expired")
	} else if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		userReq := models.UserReq{
			User:  claims.User,
			Token: tokenString,
		}
		c.Locals("user", userReq)
		return c.Next()
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Invalid token, Please try again later")
}
