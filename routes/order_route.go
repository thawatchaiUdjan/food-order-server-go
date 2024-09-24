package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderRoute(app *fiber.App, db *mongo.Database) {
	orderService := services.CreateOrderService(db)
	route := app.Group("/orders", middlewares.AuthToken)

	route.Get("/", orderService.FindOne)
	route.Get("/all-order", orderService.FindAll)
	route.Post("/", orderService.Create)
	route.Put(":id/:status", orderService.UpdateStatus)
	route.Delete("/:id", orderService.Remove)
}
