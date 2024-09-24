package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderStatusRoute(app *fiber.App, db *mongo.Database) {
	orderStatusService := services.CreateOrderStatusService(db)
	route := app.Group("/order-status", middlewares.AuthToken)

	route.Get("/", orderStatusService.FindAll)
}
