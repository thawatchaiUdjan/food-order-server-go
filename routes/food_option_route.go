package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodOptionRoute(app *fiber.App, db *mongo.Database) {
	foodOptionService := services.CreateFoodOptionService(db)
	route := app.Group("/food-option", middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		foodOptions, err := foodOptionService.FindAll()
		if err != nil {
			return fiber.NewError(fiber.StatusNotAcceptable, err.Error())
		}
		return c.JSON(foodOptions)
	})
}
