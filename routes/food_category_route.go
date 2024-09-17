package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodCategoryRoute(app *fiber.App, db *mongo.Database) {
	foodCategoryService := services.CreateFoodCategoryService(db)
	route := app.Group("/category", middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		foodCategories, err := foodCategoryService.FindAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(foodCategories)
	})
}
