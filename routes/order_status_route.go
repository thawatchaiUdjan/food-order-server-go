package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderStatusRoute(app *fiber.App, db *mongo.Database) {
	orderStatusService := services.CreateOrderStatusService(db)
	route := app.Group("/order-status", middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		orderStatus, err := orderStatusService.FindAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(orderStatus)
	})
}
