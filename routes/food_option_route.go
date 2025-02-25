package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodOptionRoute(app *fiber.App, db *mongo.Database) {
	foodOptionService := services.CreateFoodOptionService(db)
	route := app.Group("/food-option", middlewares.AuthToken)

	route.Get("/", foodOptionService.FindAll)
}
