package routes

import (
	"fmt"

	"github.com/food-order-server/models"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoute(app *fiber.App, db *mongo.Database) {
	userService := services.CreateUserService(db)
	userRoute := app.Group("/user")

	userRoute.Post("/login", func(c fiber.Ctx) error {
		userBody := new(models.UserLoginReq)

		if err := c.Bind().Body(userBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.MessageRes{
				Message: "An error occurred. Invalid request body",
			})
		}

		user, err := userService.Login(userBody)
		if err == fiber.ErrBadRequest {
			return c.Status(fiber.StatusBadRequest).JSON(models.MessageRes{
				Message: "Username or password invalid",
			})
		} else if err != nil {
			fmt.Print(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(models.MessageRes{
				Message: "An unexpected error occurred. Please try again later",
			})
		}

		return c.JSON(user)
	})

	userRoute.Post("/register", func(c fiber.Ctx) error {
		userBody := new(models.UserRegisterReq)

		if err := c.Bind().Body(userBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.MessageRes{
				Message: "An error occurred. Invalid request body",
			})
		}

		user, err := userService.Register(userBody)
		if err == fiber.ErrConflict {
			return c.Status(fiber.StatusBadRequest).JSON(models.MessageRes{
				Message: "Username is already in use",
			})
		} else if err != nil {
			fmt.Print(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(models.MessageRes{
				Message: "An unexpected error occurred. Please try again later",
			})
		}

		return c.JSON(user)
	})
}
