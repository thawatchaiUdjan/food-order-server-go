package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodRoute(app *fiber.App, db *mongo.Database) {
	foodService := services.CreateFoodService(db)
	route := app.Group("/foods", middlewares.AuthToken)

	route.Get("/", foodService.FindAll)
	route.Post("/", middlewares.UploadFoodFile, foodService.Create)
	route.Put("/:id", middlewares.UploadFoodFile, foodService.Update)
	route.Delete("/:id", foodService.Remove)
}
