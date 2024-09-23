package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate = validator.New()

func Validate(data interface{}) error {
	if err := validate.Struct(data); err != nil {
		return fiber.ErrBadRequest
	}
	return nil
}
