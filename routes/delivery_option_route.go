package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeliveryOptionRoute(app *fiber.App, db *mongo.Database) {
	deliveryOptionService := services.CreateDeliveryOptionService(db)
	route := app.Group("/delivery", middlewares.AuthToken)

	route.Get("/", deliveryOptionService.FindAll)
}
