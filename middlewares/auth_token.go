package middlewares

import (
	"strings"

	"github.com/food-order-server/config"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthToken(c fiber.Ctx) error {
	config := config.LoadConfig()
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No token provided. Access denied")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err == jwt.ErrTokenExpired {
		return fiber.NewError(fiber.StatusUnauthorized, "Token has expired")
	} else if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("user", claims)
		return c.Next()
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
}
