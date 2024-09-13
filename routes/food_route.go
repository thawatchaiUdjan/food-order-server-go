package routes

import (
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodRoute(app *fiber.App, db *mongo.Database) {
	foodService := services.CreateFoodService(db)
	foodRoute := app.Group("/foods")

	foodRoute.Get("/", func(c fiber.Ctx) error {
		foods, err := foodService.GetFoods()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(foods)
	})
}
