package routes

import (
	"github.com/food-order-server/models"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodRoute(app *fiber.App, db *mongo.Client) {
	foodService := services.CreateFoodService(db)

	app.Get("/foods", func(c fiber.Ctx) error {
		foods, err := foodService.GetFoods()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MessageRes{Message: "An unexpected error occurred. Please try again later" + err.Error()})
		}
		return c.JSON(foods)
	})
}
