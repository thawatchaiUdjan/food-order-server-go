package middlewares

import (
	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	var status int
	var message string

	if e, ok := err.(*fiber.Error); ok {
		status = e.Code
		message = e.Message
	}

	switch status {
	case fiber.StatusBadRequest:
		message = "An error occurred. Invalid request body."
	case fiber.StatusInternalServerError:
		message = "An unexpected error occurred. Please try again later."
	}

	return c.Status(status).JSON(models.MessageRes{
		Message: message,
	})
}
