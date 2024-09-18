package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeliveryOptionRoute(app *fiber.App, db *mongo.Database) {
	deliveryOptionService := services.CreateDeliveryOptionService(db)
	route := app.Group("/delivery", middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		deliveryOptions, err := deliveryOptionService.FindAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(deliveryOptions)
	})
}
