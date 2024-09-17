package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodRoute(app *fiber.App, db *mongo.Database) {
	foodService := services.CreateFoodService(db)
	route := app.Group("/foods", middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		foods, err := foodService.FindAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(foods)
	})
}
